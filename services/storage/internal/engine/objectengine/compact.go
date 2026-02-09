package objectengine

import (
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log/slog"
	"os"
	"time"
)

func (b *Bucket) startCompactSchedule() {
	go func() {
		for {
			slog.Info("Starting scheduled compaction for bucket", slog.String("database", b.database), slog.String("collection", b.collection))

			time.Sleep(10 * time.Second)

			for i, s := range b.shards {
				if !s.dirty {
					continue
				}

				slog.Debug("Compacting shard", slog.Int("shard_index", i))
				if err := b.compactShard(s); err != nil {
					slog.Error("Failed to compact shard", slog.Int("shard_index", i), slog.Any("error", err))
				}

				time.Sleep(10 * time.Second) // Sleep between shard compactions to reduce load
			}

			time.Sleep(5 * time.Minute)
		}
	}()
}

func (b *Bucket) compactShard(s *shard) error {
	s.Lock()
	defer s.Unlock()

	if !s.dirty {
		return nil
	}

	// Rewind old file
	if _, err := s.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	shardPath := s.file.Name()
	tmpPath := s.file.Name() + ".tmp"
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	newIndex := make(map[string]ObjectMeta)

	for {
		entryOffset, _ := s.file.Seek(0, io.SeekCurrent)

		meta, _, err := readEntry(s.file, true)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Only keep the version with the latest offset (most recent) for each key
		current, ok := s.index[meta.Key]
		if ok {
			// If the current index entry has a higher offset than the one we're reading, it means the current entry is more recent and we should skip the older one.
			if current.Offset > entryOffset {
				slog.Debug("Skipping older entry during compaction", slog.String("key", meta.Key), slog.Int64("current_offset", current.Offset), slog.Int64("entry_offset", entryOffset))
				continue
			}
		}

		newIndex[meta.Key] = ObjectMeta{
			Key:        meta.Key,
			Offset:     entryOffset,
			Size:       meta.Size,
			Checksum:   meta.Checksum,
			CreatedAt:  meta.CreatedAt,
			ModifiedAt: meta.ModifiedAt,
		}
	}

	// Write all valid entries to the new shard file
	for _, meta := range newIndex {
		fmt.Printf("Compacting entry: key=%s, offset=%d, size=%d\n", meta.Key, meta.Offset, meta.Size)

		// Seek to the offset and read the entry
		if _, err := s.file.Seek(meta.Offset, io.SeekStart); err != nil {
			return err
		}

		// Read the entry and verify the key matches
		readMeta, data, err := readEntry(s.file, false)
		if err != nil {
			return err
		}
		if readMeta.Key != meta.Key {
			return errors.New("key mismatch")
		}

		// Verify checksum
		if crc32.ChecksumIEEE(data) != meta.Checksum {
			return errors.New("data checksum mismatch")
		}

		if err := writeEntry(tmpFile, readMeta.Key, readMeta.Checksum, readMeta.CreatedAt, readMeta.ModifiedAt, data); err != nil {
			return err
		}
	}

	// Ensure data is durable
	if err := tmpFile.Sync(); err != nil {
		return err
	}

	// Close old shard before replacing
	oldFile := s.file
	if err := oldFile.Close(); err != nil {
		return err
	}

	// Atomic replace
	if err := os.Rename(tmpPath, shardPath); err != nil {
		return err
	}

	// Reopen shard file
	f, err := os.OpenFile(shardPath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	s.file = f
	s.index = newIndex
	s.dirty = false

	return nil
}
