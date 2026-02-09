package objectengine

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fancyinnovations/fancyspaces/storage/internal/hashing"
)

// ShardCount defines the number of shards for partitioning data
const ShardCount = 16

// shard holds a file and an in-memory index
type shard struct {
	sync.RWMutex
	file  *os.File
	index map[string]ObjectMeta
	dirty bool // Indicates if the shard has deleted entries and needs compaction
}

// shardForKey determines which shard a key belongs to using FNV-1a hash
func (b *Bucket) shardForKey(key string) *shard {
	h := hashing.FNV32a(key)
	return b.shards[int(h%ShardCount)]
}

// newShard creates or opens a shard file and loads its index
// Shard files are named shard_0.bin, shard_1.bin, ..., shard_15.bin
func (b *Bucket) newShard(idx int, basePath string) (*shard, error) {
	shardFilename := fmt.Sprintf("shard_%d.bin", idx)
	shardPath := filepath.Join(basePath, shardFilename)

	f, err := os.OpenFile(shardPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	s := &shard{
		file:  f,
		index: make(map[string]ObjectMeta),
		dirty: true, // Mark as dirty to trigger initial compaction if there are deleted entries
	}

	if err := s.loadIndex(); err != nil {
		return nil, err
	}

	return s, nil
}

// loadIndex reads all entries from the shard file and builds the in-memory index.
// It also marks the shard as dirty if it encounters any deleted entries (size == 0).
func (s *shard) loadIndex() error {
	idx := make(map[string]ObjectMeta)

	entries, err := readAllEntries(s)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.Size == 0 { // Deleted entry, skip but mark shard as dirty
			s.dirty = true
			continue
		}

		existing, exists := idx[e.Key]
		if exists {
			// If we encounter multiple entries for the same key, we keep the one with the highest offset (most recent)
			if existing.Offset > e.Offset {
				continue
			}
		}

		idx[e.Key] = *e
	}

	// Swap in the rebuilt index
	s.index = idx

	return nil
}
