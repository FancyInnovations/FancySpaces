package broker

import (
	"sync"
)

// queueGroupCounter is a global counter for assigning unique IDs to queue subscribers
var queueGroupCounter uint32

// Subscriber represents a subscriber
type Subscriber struct {
	ID    string
	Queue string
	msgCh chan []byte
}

// Node represents a trie node for subject routing
type Node struct {
	sync.RWMutex
	subs     []*Subscriber
	children map[string]*Node
	star     *Node // *
	greater  *Node // >
}
