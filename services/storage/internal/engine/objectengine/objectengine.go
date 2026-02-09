package objectengine

import (
	"errors"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
)

// Bucket represents a persistent bucket with sharded storage
type Bucket struct {
	database   string
	collection string
	shards     [ShardCount]*shard
}

type Configuration struct {
	Database   string
	Collection string
}

// NewBucket initializes a new bucket, creating necessary directories and shard files
// Path structure: data/{database}/{collection}/shard_{i}.bin
func NewBucket(cfg Configuration) (*Bucket, error) {
	b := &Bucket{
		database:   cfg.Database,
		collection: cfg.Collection,
	}
	path := filepath.Join("data", b.database, b.collection)

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	// Initialize shards
	for i := 0; i < ShardCount; i++ {
		s, err := b.newShard(i, path)
		if err != nil {
			return nil, err
		}

		b.shards[i] = s
	}

	// Start background compaction
	b.startCompactSchedule()

	return b, nil
}

// Put stores an object in the bucket
func (b *Bucket) Put(key string, data []byte) error {
	s := b.shardForKey(key)
	s.Lock()
	defer s.Unlock()

	offset, err := s.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	// Write length-prefixed key + value
	if err := writeEntry(s.file, key, data); err != nil {
		return err
	}

	// Compute checksum
	checksum := crc32.ChecksumIEEE(data)

	// Update in-memory index
	s.index[key] = ObjectMeta{
		Offset:   offset,
		Size:     int64(len(data)),
		Checksum: checksum,
	}

	return nil
}

// Get retrieves an object by key
func (b *Bucket) Get(key string) ([]byte, error) {
	s := b.shardForKey(key)
	s.RLock()
	defer s.RUnlock()

	meta, ok := s.index[key]
	if !ok {
		return nil, ErrKeyNotFound
	}

	// Seek to the offset and read the entry
	if _, err := s.file.Seek(meta.Offset, io.SeekStart); err != nil {
		return nil, err
	}

	// Read the entry and verify the key matches
	readKey, data, err := readEntry(s.file)
	if err != nil {
		return nil, err
	}
	if readKey != key {
		return nil, errors.New("key mismatch")
	}

	// Verify checksum
	if crc32.ChecksumIEEE(data) != meta.Checksum {
		return nil, errors.New("data checksum mismatch")
	}

	return data, nil
}

// Delete marks a key as deleted
func (b *Bucket) Delete(key string) error {
	s := b.shardForKey(key)
	s.Lock()
	defer s.Unlock()

	if _, ok := s.index[key]; !ok {
		return ErrKeyNotFound
	}

	delete(s.index, key)

	s.dirty = true // Mark shard as dirty for compaction

	return nil
}

// GetMeta returns the ObjectMeta for a given key
func (b *Bucket) GetMeta(key string) (*ObjectMeta, error) {
	s := b.shardForKey(key)
	s.RLock()
	defer s.RUnlock()

	meta, ok := s.index[key]
	if !ok {
		return nil, ErrKeyNotFound
	}

	return &meta, nil
}
