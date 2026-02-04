package kv

import "sync"

const ShardCount = 256

type entry struct {
	value   Value
	expires int64 // unix nanos, 0 means no expiration
}

type shard struct {
	index int
	mu    sync.RWMutex
	data  map[string]*entry
}
