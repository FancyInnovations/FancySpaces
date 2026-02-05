package kv

import (
	"testing"
	"time"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
)

func contains(slice []string, v string) bool {
	for _, s := range slice {
		if s == v {
			return true
		}
	}
	return false
}

func TestSetGetExistsDeleteClearSizeKeys(t *testing.T) {
	e := NewEngine(Configuration{
		DisableTTL: false,
	})

	var v codex.Value

	// initially empty
	if e.Size() != 0 {
		t.Fatalf("expected size 0, got %d", e.Size())
	}
	if len(e.Keys()) != 0 {
		t.Fatalf("expected no keys, got %v", e.Keys())
	}

	// Set and verify
	e.Set("key1", &v)
	if !e.Exists("key1") {
		t.Fatalf("expected key1 to exist")
	}
	if got := e.Get("key1"); got == nil {
		t.Fatalf("expected Get(key1) to be non-nil")
	}

	if e.Size() != 1 {
		t.Fatalf("expected size 1, got %d", e.Size())
	}
	keys := e.Keys()
	if !contains(keys, "key1") {
		t.Fatalf("expected keys to contain key1, got %v", keys)
	}

	// Delete
	e.Delete("key1")
	if e.Exists("key1") {
		t.Fatalf("expected key1 to not exist after delete")
	}
	if e.Get("key1") != nil {
		t.Fatalf("expected Get(key1) == nil after delete")
	}
	if e.Size() != 0 {
		t.Fatalf("expected size 0 after delete, got %d", e.Size())
	}

	// Clear
	e.Set("a", &v)
	e.Set("b", &v)
	if e.Size() != 2 {
		t.Fatalf("expected size 2, got %d", e.Size())
	}
	e.Clear()
	if e.Size() != 0 {
		t.Fatalf("expected size 0 after clear, got %d", e.Size())
	}
	if len(e.Keys()) != 0 {
		t.Fatalf("expected no keys after clear, got %v", e.Keys())
	}
}

func TestGetMultipleAndDeleteMultiple(t *testing.T) {
	e := NewEngine(Configuration{
		DisableTTL: false,
	})
	var v codex.Value

	e.Set("k1", &v)
	e.Set("k2", &v)
	e.Set("k3", &v)

	res := e.GetMultiple([]string{"k1", "k2", "missing"})
	if res["k1"] == nil || res["k2"] == nil {
		t.Fatalf("expected k1 and k2 to be present in GetMultiple")
	}
	if _, ok := res["missing"]; !ok {
		t.Fatalf("expected result to contain key for missing entry")
	}
	if res["missing"].Type != codex.TypeEmpty {
		t.Fatalf("expected missing entry to be nil")
	}

	e.DeleteMultiple([]string{"k1", "k3"})
	if e.Exists("k1") || e.Exists("k3") {
		t.Fatalf("expected k1 and k3 to be removed")
	}
	// k2 should remain
	if !e.Exists("k2") {
		t.Fatalf("expected k2 to remain")
	}
}

func TestSetIfExistsSetIfNotExistsWithExpiry(t *testing.T) {
	e := NewEngine(Configuration{
		DisableTTL: false,
	})
	var v codex.Value

	// Not exists -> SetIfExistsTTL should be false, SetIfNotExistsTTL true
	if e.SetIfExistsTTL("x", v, 0) {
		t.Fatalf("SetIfExistsTTL should be false for non-existing key")
	}
	if !e.SetIfNotExistsTTL("x", v, 0) {
		t.Fatalf("SetIfNotExistsTTL should be true for non-existing key")
	}
	if !e.SetIfExistsTTL("x", v, 0) {
		t.Fatalf("SetIfExistsTTL should be true for existing key")
	}

	// Expiry: set with past expiration
	past := time.Now().Add(-1 * time.Second).UnixNano()
	e.SetWithTTL("exp", &v, past)
	if e.Exists("exp") {
		t.Fatalf("expected exp to be treated as expired by Exists")
	}
	// expired -> SetIfNotExistsTTL should succeed (treat as not exists)
	if !e.SetIfNotExistsTTL("exp", v, 0) {
		t.Fatalf("expected SetIfNotExistsTTL to set expired key")
	}
	// now exists
	if !e.SetIfExistsTTL("exp", v, 0) {
		t.Fatalf("expected SetIfExistsTTL to succeed for unexpired key")
	}
}
