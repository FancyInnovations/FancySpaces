package kv

import "time"

func (e *Engine) cleanup() {
	for i := 0; i < ShardCount; i++ {
		s := &e.shards[i]
		s.mu.Lock()

		now := time.Now().UnixNano()

		for key, entry := range s.data {
			if entry.expires > 0 && entry.expires <= now {
				delete(s.data, key)
			}
		}
		s.mu.Unlock()
	}
}

func (e *Engine) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			e.cleanup()
		}
	}()
}
