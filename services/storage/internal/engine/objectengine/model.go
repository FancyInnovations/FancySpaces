package objectengine

// ObjectMeta stores object metadata including checksum
type ObjectMeta struct {
	Offset   int64
	Size     int64
	Checksum uint32 // CRC32 checksum of the data

	// TODO: Add more metadata fields if needed (e.g., creation time, last modified time, etc.)
}
