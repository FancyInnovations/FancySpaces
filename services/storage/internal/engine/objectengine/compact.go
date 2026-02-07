package objectengine

import (
	"io"
	"log/slog"
	"os"
	"time"
)

func (b *Bucket) startCompactSchedule() {
	go func() {
		for {
			slog.Info("Starting scheduled compaction for bucket", slog.String("database", b.database), slog.String("collection", b.collection))

			time.Sleep(5 * time.Minute)

			for i, s := range b.shards {
				if !s.dirty {
					continue
				}

				if err := b.Compact(s); err != nil {
					slog.Error("Failed to compact shard", slog.Int("shard_index", i), slog.Any("error", err))
				}

				time.Sleep(10 * time.Second) // Sleep between shard compactions to reduce load
			}
		}
	}()
}

// Compact rewrites the shard file to remove deleted entries and reduce fragmentation
func (b *Bucket) Compact(s *shard) error {
	s.Lock()

	// Create temp file
	tmpPath := s.file.Name() + ".tmp"
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		s.Unlock()
		return err
	}

	// Iterate live objects in index and write them
	for key, meta := range s.index {
		if _, err := s.file.Seek(meta.Offset, io.SeekStart); err != nil {
			tmpFile.Close()
			s.Unlock()
			return err
		}
		_, data, err := readEntry(s.file)
		if err != nil {
			tmpFile.Close()
			s.Unlock()
			return err
		}
		if err := writeEntry(tmpFile, key, data); err != nil {
			tmpFile.Close()
			s.Unlock()
			return err
		}

		// Update meta offset in new file
		offset, _ := tmpFile.Seek(0, io.SeekCurrent)
		offset -= int64(4 + len(key) + 4 + len(data)) // key length (4 bytes) + key + data length (4 bytes) + data
		s.index[key] = ObjectMeta{
			Offset:   offset,
			Size:     meta.Size,
			Checksum: meta.Checksum,
		}
	}

	tmpFile.Sync()
	tmpFile.Close()
	s.file.Close()

	// Replace old shard file atomically
	if err := os.Rename(tmpPath, s.file.Name()); err != nil {
		s.Unlock()
		return err
	}

	// Reopen the new file
	s.file, err = os.OpenFile(s.file.Name(), os.O_RDWR, 0644)
	if err != nil {
		s.Unlock()
		return err
	}

	// Mark shard as clean
	s.dirty = false

	s.Unlock()
	return nil
}
