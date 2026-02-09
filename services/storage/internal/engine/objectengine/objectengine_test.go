package objectengine

import (
	"bytes"
	"hash/crc32"
	"os"
	"testing"
)

func TestPutGetGetMetaDelete_TableDriven(t *testing.T) {
	cfg := Configuration{
		Database:   "testdb",
		Collection: "testcol",
	}

	b, err := NewBucket(cfg)
	if err != nil {
		t.Fatalf("NewBucket error: %v", err)
	}

	cases := []struct {
		name  string
		key   string
		value []byte
	}{
		{name: "small", key: "k1", value: []byte("v1")},
		{name: "medium", key: "k2", value: []byte("some longer value")},
		{name: "empty", key: "k3", value: []byte("")},
	}

	// Put all entries
	for _, tc := range cases {
		tc := tc
		t.Run("Put_"+tc.name, func(t *testing.T) {
			if err := b.Put(tc.key, tc.value); err != nil {
				t.Fatalf("Put(%s) error: %v", tc.key, err)
			}
		})
	}

	// Get and GetMeta for all entries
	for _, tc := range cases {
		tc := tc
		t.Run("Get_"+tc.name, func(t *testing.T) {
			got, err := b.Get(tc.key)
			if err != nil {
				t.Fatalf("Get(%s) error: %v", tc.key, err)
			}
			if !bytes.Equal(got, tc.value) {
				t.Fatalf("Get(%s) returned wrong value: got=%v want=%v", tc.key, got, tc.value)
			}

			meta, err := b.GetMeta(tc.key)
			if err != nil {
				t.Fatalf("GetMeta(%s) error: %v", tc.key, err)
			}
			if meta.Size != uint32(len(tc.value)) {
				t.Fatalf("GetMeta(%s) size mismatch: got=%d want=%d", tc.key, meta.Size, len(tc.value))
			}
			if meta.Checksum != crc32.ChecksumIEEE(tc.value) {
				t.Fatalf("GetMeta(%s) checksum mismatch: got=%d want=%d", tc.key, meta.Checksum, crc32.ChecksumIEEE(tc.value))
			}
		})
	}

	// Delete a key and ensure it's gone
	t.Run("Delete_and_GetMissing", func(t *testing.T) {
		key := cases[0].key
		if err := b.Delete(key); err != nil {
			t.Fatalf("Delete(%s) error: %v", key, err)
		}

		if _, err := b.Get(key); err == nil {
			t.Fatalf("expected Get(%s) to return error after delete", key)
		}

		if _, err := b.GetMeta(key); err == nil {
			t.Fatalf("expected GetMeta(%s) to return error after delete", key)
		}
	})

	// Ensure asking for a completely missing key returns the expected error
	t.Run("Get_missing_key", func(t *testing.T) {
		_, err := b.Get("non-existent-key")
		if err == nil {
			t.Fatalf("expected error for missing key")
		}

	})

	deleteShardsForTest(t, b)
}

func deleteShardsForTest(t *testing.T, b *Bucket) {
	for _, s := range b.shards {
		file := s.file
		if err := file.Close(); err != nil {
			t.Logf("error closing shard file: %v", err)
		}
		if err := os.Remove(file.Name()); err != nil {
			t.Logf("error deleting shard file: %v", err)
		}
	}
}
