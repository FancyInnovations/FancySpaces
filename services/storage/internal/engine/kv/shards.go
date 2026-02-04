package kv

func fnv32a(key string) uint32 {
	const (
		offset uint32 = 2166136261
		prime  uint32 = 16777619
	)
	var h uint32 = offset
	for i := 0; i < len(key); i++ {
		h ^= uint32(key[i])
		h *= prime
	}
	return h
}

func (e *Engine) shardFor(key string) *shard {
	h := fnv32a(key)
	idx := int(h % uint32(ShardCount))
	return &e.shards[idx]
}
