package kv

import (
	"time"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
)

type Engine struct {
	shards     [ShardCount]shard
	disableTTL bool
}

type Configuration struct {
	DisableTTL bool
}

func NewEngine(cfg Configuration) *Engine {
	e := &Engine{
		disableTTL: cfg.DisableTTL,
	}

	for i := 0; i < ShardCount; i++ {
		e.shards[i] = shard{
			index: i,
			data:  make(map[string]*entry),
		}
	}

	if !e.disableTTL {
		e.startCleanup(1 * time.Second)
	}

	return e
}

// Keys returns a slice of all keys in the engine.
func (e *Engine) Keys() []string {
	keys := make([]string, 0)
	for i := 0; i < ShardCount; i++ {
		s := &e.shards[i]
		s.mu.RLock()
		for k := range s.data {
			keys = append(keys, k)
		}
		s.mu.RUnlock()
	}

	return keys
}

// SetWithTTL stores a value for the given key with an optional expiration time (in unix nanoseconds).
// If expires is 0, the key will not expire.
// If expires is a positive value, it should be a unix timestamp in nanoseconds indicating when the key should expire.
func (e *Engine) SetWithTTL(key string, value *codex.Value, expires int64) {
	if e.disableTTL {
		expires = 0
	}

	s := e.shardFor(key)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = &entry{
		value:   value,
		expires: expires,
	}
}

func (e *Engine) Set(key string, value *codex.Value) {
	e.SetWithTTL(key, value, 0)
}

// SetMultipleTTL allows setting multiple key-value pairs at once with the same expiration time.
// This is more efficient than calling Set multiple times, as it minimizes locking overhead by grouping updates by shard.
func (e *Engine) SetMultipleTTL(entries map[string]codex.Value, expires int64) {
	if e.disableTTL {
		expires = 0
	}

	// find shards that need to be updated
	shardEntries := make(map[int]map[string]codex.Value)
	for key, value := range entries {
		s := e.shardFor(key)
		if shardEntries[s.index] == nil {
			shardEntries[s.index] = make(map[string]codex.Value)
		}
		shardEntries[s.index][key] = value
	}

	// update each shard
	for shardIndex, entries := range shardEntries {
		s := &e.shards[shardIndex]
		s.mu.Lock()
		for key, value := range entries {
			s.data[key] = &entry{
				value:   &value,
				expires: expires,
			}
		}
		s.mu.Unlock()
	}
}

// SetMultiple allows setting multiple key-value pairs at once without expiration.
func (e *Engine) SetMultiple(entries map[string]codex.Value) {
	e.SetMultipleTTL(entries, 0)
}

// SetIfExistsTTL updates the value for the given key only if it already exists and has not expired.
// Returns true if the key was updated, false otherwise.
func (e *Engine) SetIfExistsTTL(key string, value codex.Value, expires int64) bool {
	if e.disableTTL {
		expires = 0
	}

	s := e.shardFor(key)
	s.mu.Lock()
	defer s.mu.Unlock()

	en, exists := s.data[key]
	if !exists || (en.expires > 0 && time.Now().UnixNano() > en.expires) {
		return false
	}

	s.data[key] = &entry{
		value:   &value,
		expires: expires,
	}
	return true
}

// SetIfExists updates the value for the given key only if it already exists and has not expired.
// Returns true if the key was updated, false otherwise.
func (e *Engine) SetIfExists(key string, value codex.Value) bool {
	return e.SetIfExistsTTL(key, value, 0)
}

// SetIfNotExistsTTL sets the value for the given key only if it does not already exist or has expired.
// Returns true if the key was set, false otherwise.
func (e *Engine) SetIfNotExistsTTL(key string, value codex.Value, expires int64) bool {
	if e.disableTTL {
		expires = 0
	}

	s := e.shardFor(key)
	s.mu.Lock()
	defer s.mu.Unlock()

	en, exists := s.data[key]
	if exists && (en.expires == 0 || time.Now().UnixNano() <= en.expires) {
		return false
	}

	s.data[key] = &entry{
		value:   &value,
		expires: expires,
	}
	return true
}

func (e *Engine) SetIfNotExists(key string, value codex.Value) bool {
	return e.SetIfNotExistsTTL(key, value, 0)
}

// Get retrieves the value for the given key.
func (e *Engine) Get(key string) *codex.Value {
	s := e.shardFor(key)
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.data[key]
	if !exists || (entry.expires > 0 && time.Now().UnixNano() > entry.expires) {
		return nil
	}

	return entry.value
}

// GetMultiple retrieves values for multiple keys at once.
// This is more efficient than calling Get multiple times, as it minimizes locking overhead by grouping reads by shard.
func (e *Engine) GetMultiple(keys []string) map[string]*codex.Value {
	var shardKeys [ShardCount][]string
	for _, key := range keys {
		s := e.shardFor(key)
		shardKeys[s.index] = append(shardKeys[s.index], key)
	}

	results := make(map[string]*codex.Value, len(keys))
	now := time.Now().UnixNano()
	for shardIndex, keys := range shardKeys {
		s := &e.shards[shardIndex]
		s.mu.RLock()
		for _, key := range keys {
			entry, exists := s.data[key]
			if exists && (entry.expires == 0 || now <= entry.expires) {
				results[key] = entry.value
			} else {
				results[key] = codex.EmptyValue // indicate missing/expired keys with a null value
			}
		}
		s.mu.RUnlock()
	}

	return results
}

// GetAll retrieves all key-value pairs currently stored in the engine.
// This is a heavy operation and should be used with caution, as it will read all data in memory and may block other operations while it runs.
func (e *Engine) GetAll() map[string]*codex.Value {
	results := make(map[string]*codex.Value)
	now := time.Now().UnixNano()
	for i := 0; i < ShardCount; i++ {
		s := &e.shards[i]
		s.mu.RLock()
		for key, entry := range s.data {
			if entry.expires == 0 || now <= entry.expires {
				results[key] = entry.value
			}
		}
		s.mu.RUnlock()
	}

	return results
}

// Exists checks if a key exists and has not expired.
func (e *Engine) Exists(key string) bool {
	s := e.shardFor(key)
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.data[key]
	return exists && (entry.expires == 0 || time.Now().UnixNano() <= entry.expires)
}

// Delete removes a key from the engine.
func (e *Engine) Delete(key string) {
	s := e.shardFor(key)
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}

// DeleteMultiple removes multiple keys from the engine at once.
// This is more efficient than calling Delete multiple times, as it minimizes locking overhead by grouping deletions by shard.
func (e *Engine) DeleteMultiple(keys []string) {
	shardKeys := make(map[int][]string)
	for _, key := range keys {
		s := e.shardFor(key)
		shardKeys[s.index] = append(shardKeys[s.index], key)
	}

	for shardIndex, keys := range shardKeys {
		s := &e.shards[shardIndex]
		s.mu.Lock()
		for _, key := range keys {
			delete(s.data, key)
		}
		s.mu.Unlock()
	}
}

// Clear removes all keys from the engine.
// This is a heavy operation and should be used with caution, as it will block all operations while it runs.
func (e *Engine) Clear() {
	for i := 0; i < ShardCount; i++ {
		s := &e.shards[i]
		s.mu.Lock()
		s.data = make(map[string]*entry)
		s.mu.Unlock()
	}
}

// Size returns the total number of keys currently stored in the engine across all shards.
func (e *Engine) Size() int {
	total := 0
	for i := 0; i < ShardCount; i++ {
		s := &e.shards[i]
		s.mu.RLock()
		total += len(s.data)
		s.mu.RUnlock()
	}
	return total
}
