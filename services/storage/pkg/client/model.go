package client

type ObjectMetadata struct {
	Size     int64  `json:"size"`
	Checksum uint32 `json:"checksum"`
}
