package kv

import (
	"sync"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
)

const ShardCount = 256

type entry struct {
	value   *codex.Value
	expires int64 // unix nanos, 0 means no expiration
}

type shard struct {
	index int
	mu    sync.RWMutex
	data  map[string]*entry
}
