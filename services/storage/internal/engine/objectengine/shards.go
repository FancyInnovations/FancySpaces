package objectengine

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
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

// loadIndex rebuilds the in-memory index by scanning the file
func (s *shard) loadIndex() error {
	s.index = make(map[string]ObjectMeta)

	var offset int64 = 0
	for {
		startOffset := offset

		// Read key length and key
		var keyLen uint32
		if err := binary.Read(s.file, binary.LittleEndian, &keyLen); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		keyBuf := make([]byte, keyLen)
		if _, err := io.ReadFull(s.file, keyBuf); err != nil {
			return err
		}
		key := string(keyBuf)

		// Read data length and data
		var dataLen uint32
		if err := binary.Read(s.file, binary.LittleEndian, &dataLen); err != nil {
			return err
		}

		dataBuf := make([]byte, dataLen)
		if _, err := io.ReadFull(s.file, dataBuf); err != nil {
			return err
		}

		// Compute checksum
		checksum := crc32.ChecksumIEEE(dataBuf)

		s.index[key] = ObjectMeta{
			Offset:   startOffset,
			Size:     int64(dataLen),
			Checksum: checksum,
		}

		offset, _ = s.file.Seek(0, io.SeekCurrent)
	}
	return nil
}
