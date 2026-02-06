package broker

import (
	"strings"
	"sync/atomic"
	"time"
)

type PublishCallback func(sub *Subscriber, subject string, msgs [][]byte)

type Broker struct {
	root            *Node
	pubCallback     PublishCallback
	isClientHealthy func(id string) bool
	batchSize       int
	batchTimeout    time.Duration
}

type Configuration struct {
	PublishCallback PublishCallback
	IsClientHealthy func(id string) bool
	BatchSize       int
	BatchTimeout    time.Duration
}

// NewBroker creates a new broker with a callback invoked for delivery
func NewBroker(cfg Configuration) *Broker {
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 10
	}
	if cfg.BatchTimeout <= 0 {
		cfg.BatchTimeout = 100 * time.Millisecond
	}
	if cfg.IsClientHealthy == nil {
		cfg.IsClientHealthy = func(id string) bool { return true }
	}

	b := &Broker{
		root: &Node{
			children: make(map[string]*Node),
		},
		pubCallback:     cfg.PublishCallback,
		isClientHealthy: cfg.IsClientHealthy,
		batchSize:       cfg.BatchSize,
		batchTimeout:    cfg.BatchTimeout,
	}

	b.startCleanupUnhealthyClients()

	return b
}

// Subscribe adds a subscriber to a subject and starts delivery goroutine
func (b *Broker) Subscribe(subject string, sub *Subscriber) {
	sub.msgCh = make(chan []byte, 1024)

	b.startDelivery(sub, subject, b.pubCallback)

	tokens := strings.Split(subject, ".")
	b.insert(b.root, tokens, sub)
}

// startDelivery starts the subscriber goroutine that batches messages
func (b *Broker) startDelivery(sub *Subscriber, subject string, callback PublishCallback) {
	go func() {
		batch := make([][]byte, 0, b.batchSize)
		timer := time.NewTimer(b.batchTimeout)
		defer timer.Stop()

		for {
			select {
			case msg, ok := <-sub.msgCh:
				if !ok {
					if len(batch) > 0 {
						callback(sub, subject, batch)
					}
					return
				}
				batch = append(batch, msg)
				if len(batch) >= b.batchSize {
					callback(sub, subject, batch)
					batch = batch[:0]
					timer.Reset(b.batchTimeout)
				}
			case <-timer.C:
				if len(batch) > 0 {
					callback(sub, subject, batch)
					batch = batch[:0]
				}
				timer.Reset(b.batchTimeout)
			}
		}
	}()
}

// insert recursively adds subscriber to the trie
func (b *Broker) insert(node *Node, tokens []string, sub *Subscriber) {
	if len(tokens) == 0 {
		node.Lock()
		node.subs = append(node.subs, sub)
		node.Unlock()
		return
	}

	token := tokens[0]
	var child *Node

	node.Lock()
	defer node.Unlock()

	switch token {
	case "*":
		if node.star == nil {
			node.star = &Node{children: make(map[string]*Node)}
		}
		child = node.star
	case ">":
		if node.greater == nil {
			node.greater = &Node{children: make(map[string]*Node)}
		}
		child = node.greater
	default:
		if node.children[token] == nil {
			node.children[token] = &Node{children: make(map[string]*Node)}
		}
		child = node.children[token]
	}

	b.insert(child, tokens[1:], sub)
}

// Unsubscribe removes a subscriber from a subject
func (b *Broker) Unsubscribe(subject string, subID string) {
	tokens := strings.Split(subject, ".")
	b.remove(b.root, tokens, subID)
}

// remove recursively removes subscriber
func (b *Broker) remove(node *Node, tokens []string, subID string) {
	if node == nil {
		return
	}

	if len(tokens) == 0 {
		node.Lock()
		filtered := node.subs[:0]
		for _, s := range node.subs {
			if s.ID != subID {
				filtered = append(filtered, s)
			} else {
				close(s.msgCh) // stop delivery
			}
		}
		node.subs = filtered
		node.Unlock()
		return
	}

	token := tokens[0]
	node.RLock()
	defer node.RUnlock()

	switch token {
	case "*":
		b.remove(node.star, tokens[1:], subID)
	case ">":
		b.remove(node.greater, tokens[1:], subID)
	default:
		b.remove(node.children[token], tokens[1:], subID)
	}
}

// Publish sends a message to all matching subscribers
func (b *Broker) Publish(subject string, msg []byte) {
	tokens := strings.Split(subject, ".")
	b.publish(b.root, tokens, subject, msg)
}

// publish walks the trie and delivers messages
func (b *Broker) publish(node *Node, tokens []string, subject string, msg []byte) {
	if node == nil {
		return
	}

	// Collect subscribers under lock, then release
	node.RLock()
	subs := append([]*Subscriber(nil), node.subs...) // shallow copy
	node.RUnlock()

	// Map queue group -> subscribers
	groupMap := map[string][]*Subscriber{}
	normalSubs := []*Subscriber{}
	for _, sub := range subs {
		if sub.Queue != "" {
			groupMap[sub.Queue] = append(groupMap[sub.Queue], sub)
		} else {
			normalSubs = append(normalSubs, sub)
		}
	}

	// Deliver one message per queue group
	for _, group := range groupMap {

		// try to find a healthy subscriber in the group
		for {
			idx := atomic.AddUint32(&queueGroupCounter, 1)
			sub := group[idx%uint32(len(group))]

			if b.isClientHealthy(sub.ID) {
				select {
				case sub.msgCh <- msg:
				default: // drop if full
				}
				break
			}

			// if no healthy subscriber found after one full loop, break to avoid infinite loop
			if idx%uint32(len(group)) == 0 {
				break
			}
		}
	}

	// Deliver to normal subscribers
	for _, sub := range normalSubs {
		if !b.isClientHealthy(sub.ID) {
			continue
		}

		select {
		case sub.msgCh <- msg:
		default:
		}
	}

	// recursion for next token
	if len(tokens) == 0 {
		if node.greater != nil {
			b.publish(node.greater, nil, subject, msg)
		}
		return
	}

	token := tokens[0]

	// exact match
	if child, ok := node.children[token]; ok {
		b.publish(child, tokens[1:], subject, msg)
	}
	// single-level wildcard *
	if node.star != nil {
		b.publish(node.star, tokens[1:], subject, msg)
	}
	// tail wildcard >
	if node.greater != nil {
		b.publish(node.greater, nil, subject, msg)
	}
}

func (b *Broker) startCleanupUnhealthyClients() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			b.cleanupNode(b.root)
		}
	}()
}

func (b *Broker) cleanupNode(node *Node) {
	if node == nil {
		return
	}

	filtered := node.subs[:0]
	for _, sub := range node.subs {
		if b.isClientHealthy(sub.ID) {
			filtered = append(filtered, sub)
		} else {
			close(sub.msgCh)
		}
	}
	node.subs = filtered

	for _, child := range node.children {
		b.cleanupNode(child)
	}
	if node.star != nil {
		b.cleanupNode(node.star)
	}
	if node.greater != nil {
		b.cleanupNode(node.greater)
	}
}
