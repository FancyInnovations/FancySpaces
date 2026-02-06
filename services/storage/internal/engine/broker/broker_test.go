package broker

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBasicPublishSubscribe(t *testing.T) {
	received := make(chan [][]byte, 1)

	b := NewBroker(Configuration{
		BatchSize: 1,
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			received <- msgs
		},
	})

	sub := &Subscriber{ID: "1"}
	b.Subscribe("foo.bar", sub)

	b.Publish("foo.bar", []byte("hello"))

	select {
	case msgs := <-received:
		if len(msgs) != 1 || string(msgs[0]) != "hello" {
			t.Fatalf("unexpected msgs: %v", msgs)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout waiting for message")
	}
}

func TestBatchBySize(t *testing.T) {
	received := make(chan [][]byte, 1)

	b := NewBroker(Configuration{
		BatchSize: 3,
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			received <- msgs
		},
	})

	sub := &Subscriber{ID: "1"}
	b.Subscribe("a.b", sub)

	b.Publish("a.b", []byte("1"))
	b.Publish("a.b", []byte("2"))
	b.Publish("a.b", []byte("3"))

	select {
	case msgs := <-received:
		if len(msgs) != 3 {
			t.Fatalf("expected 3 msgs, got %d", len(msgs))
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout waiting for batch")
	}
}

func TestBatchByTimeout(t *testing.T) {
	received := make(chan [][]byte, 1)

	b := NewBroker(Configuration{
		BatchSize:    10,
		BatchTimeout: 50 * time.Millisecond,
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			received <- msgs
		},
	})

	sub := &Subscriber{ID: "1"}
	b.Subscribe("x.y", sub)

	b.Publish("x.y", []byte("hello"))

	select {
	case msgs := <-received:
		if len(msgs) != 1 {
			t.Fatalf("expected 1 msg, got %d", len(msgs))
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout waiting for timeout batch")
	}
}

func TestWildcardStar(t *testing.T) {
	received := make(chan string, 1)

	b := NewBroker(Configuration{
		BatchSize: 1,
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			received <- string(msgs[0])
		},
	})

	sub := &Subscriber{ID: "1"}
	b.Subscribe("foo.*.baz", sub)

	b.Publish("foo.bar.baz", []byte("ok"))

	select {
	case msg := <-received:
		if msg != "ok" {
			t.Fatalf("unexpected msg: %s", msg)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestWildcardGreater(t *testing.T) {
	received := make(chan string, 1)

	b := NewBroker(Configuration{
		BatchSize: 1,
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			received <- string(msgs[0])
		},
	})

	sub := &Subscriber{ID: "1"}
	b.Subscribe("foo.>", sub)

	b.Publish("foo.a.b.c", []byte("hit"))

	select {
	case msg := <-received:
		if msg != "hit" {
			t.Fatalf("unexpected msg: %s", msg)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestQueueGroupSingleDelivery(t *testing.T) {
	var count1, count2 uint32

	b := NewBroker(Configuration{
		BatchSize: 1,
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			if sub.ID == "1" {
				atomic.AddUint32(&count1, 1)
			}
			if sub.ID == "2" {
				atomic.AddUint32(&count2, 1)
			}
		},
	})

	sub1 := &Subscriber{ID: "1", Queue: "workers"}
	sub2 := &Subscriber{ID: "2", Queue: "workers"}

	b.Subscribe("jobs.run", sub1)
	b.Subscribe("jobs.run", sub2)

	for i := 0; i < 100; i++ {
		b.Publish("jobs.run", []byte("job"))
	}

	time.Sleep(200 * time.Millisecond)

	total := atomic.LoadUint32(&count1) + atomic.LoadUint32(&count2)
	if total != 100 {
		t.Fatalf("expected 100 total deliveries, got %d", total)
	}
}

func TestUnsubscribeStopsDelivery(t *testing.T) {
	received := make(chan struct{}, 1)

	b := NewBroker(Configuration{
		BatchSize: 1,
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			received <- struct{}{}
		},
	})

	sub := &Subscriber{ID: "1"}
	b.Subscribe("a.b", sub)

	b.Unsubscribe("a.b", "1")
	b.Publish("a.b", []byte("test"))

	select {
	case <-received:
		t.Fatal("received message after unsubscribe")
	case <-time.After(100 * time.Millisecond):
		// success
	}
}

func TestMultipleSubscribersAllReceive(t *testing.T) {
	var mu sync.Mutex
	count := 0

	b := NewBroker(Configuration{
		BatchSize: 1,
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			mu.Lock()
			count++
			mu.Unlock()
		},
	})

	sub1 := &Subscriber{ID: "1"}
	sub2 := &Subscriber{ID: "2"}

	b.Subscribe("foo.bar", sub1)
	b.Subscribe("foo.bar", sub2)

	b.Publish("foo.bar", []byte("hello"))

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 2 {
		t.Fatalf("expected 2 deliveries, got %d", count)
	}
}
