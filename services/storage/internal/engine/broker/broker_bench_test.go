package broker

import (
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func benchmarkBroker(
	br *Broker,
	subject string,
	batchSize int,
	b *testing.B,
) {
	msg := []byte("hello")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		br.Publish(subject, msg)
	}
}

func BenchmarkPublishSingleSubscriber(b *testing.B) {
	var delivered uint64

	br := NewBroker(Configuration{
		BatchSize: 1,
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			atomic.AddUint64(&delivered, uint64(len(msgs)))
		},
	})

	sub := &Subscriber{ID: "1"}
	br.Subscribe("foo.bar", sub)

	benchmarkBroker(br, "foo.bar", 1, b)
}

func BenchmarkPublishFanout10(b *testing.B) {
	var delivered uint64

	br := NewBroker(Configuration{
		BatchSize: 1,
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			atomic.AddUint64(&delivered, uint64(len(msgs)))
		},
	})

	for i := 0; i < 10; i++ {
		br.Subscribe("foo.bar", &Subscriber{ID: strconv.Itoa(i)})
	}

	benchmarkBroker(br, "foo.bar", 1, b)
}

func BenchmarkPublishQueueGroup10(b *testing.B) {
	var delivered uint64

	br := NewBroker(Configuration{
		BatchSize: 1,
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			atomic.AddUint64(&delivered, uint64(len(msgs)))
		},
	})

	for i := 0; i < 10; i++ {
		br.Subscribe("jobs.run", &Subscriber{
			ID:    strconv.Itoa(i),
			Queue: "workers",
		})
	}

	benchmarkBroker(br, "jobs.run", 1, b)
}

func BenchmarkPublishWildcardStar(b *testing.B) {
	var delivered uint64

	br := NewBroker(Configuration{
		BatchSize: 1,
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			atomic.AddUint64(&delivered, uint64(len(msgs)))
		},
	})

	br.Subscribe("foo.*.baz", &Subscriber{ID: "1"})

	benchmarkBroker(br, "foo.bar.baz", 1, b)
}

func BenchmarkPublishWildcardGreater(b *testing.B) {
	var delivered uint64

	br := NewBroker(Configuration{
		BatchSize: 1,
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			atomic.AddUint64(&delivered, uint64(len(msgs)))
		},
	})

	br.Subscribe("foo.>", &Subscriber{ID: "1"})

	benchmarkBroker(br, "foo.a.b.c.d.e", 1, b)
}

func BenchmarkPublishBatching(b *testing.B) {
	var delivered uint64

	br := NewBroker(Configuration{
		BatchSize:    32,
		BatchTimeout: time.Second, // force size-based batching
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			atomic.AddUint64(&delivered, uint64(len(msgs)))
		},
	})

	br.Subscribe("batch.test", &Subscriber{ID: "1"})

	msg := []byte("hello")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		br.Publish("batch.test", msg)
	}

	b.StopTimer()

	// ensure final batch flush
	time.Sleep(10 * time.Millisecond)
}

func BenchmarkPublishParallel(b *testing.B) {
	var delivered uint64

	br := NewBroker(Configuration{
		BatchSize: 1,
		PublishCallback: func(sub *Subscriber, subject string, msgs [][]byte) {
			atomic.AddUint64(&delivered, uint64(len(msgs)))
		},
	})

	br.Subscribe("foo.bar", &Subscriber{ID: "1"})

	msg := []byte("hello")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			br.Publish("foo.bar", msg)
		}
	})
}
