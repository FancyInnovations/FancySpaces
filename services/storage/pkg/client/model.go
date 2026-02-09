package client

type ObjectMetadata struct {
	// Key is the unique identifier for the object
	Key string

	// Size is the length of the object data in bytes
	Size uint32

	// CRC32 checksum of the object data for integrity verification
	Checksum uint32

	// CreatedAt is the timestamp (in unix milliseconds) when the object was created
	CreatedAt int64

	// ModifiedAt is the timestamp (in unix milliseconds) when the object was last modified
	ModifiedAt int64
}
