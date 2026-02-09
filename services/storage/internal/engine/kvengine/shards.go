package kvengine

import "github.com/fancyinnovations/fancyspaces/storage/internal/hashing"

func (e *Engine) shardFor(key string) *shard {
	h := hashing.FNV32a(key)
	return &e.shards[int(h%uint32(ShardCount))]
}
