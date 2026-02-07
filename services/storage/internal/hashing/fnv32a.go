package hashing

// FNV32a implements the FNV-1a hash algorithm for 32-bit hashes.
func FNV32a(key string) uint32 {
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
