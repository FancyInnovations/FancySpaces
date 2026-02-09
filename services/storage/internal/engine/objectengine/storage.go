package objectengine

import (
	"encoding/binary"
	"io"
)

// writeEntry writes a key-value pair with length prefixes
// Format: | key length (4 bytes) | key (variable) | checksum (4 bytes) | created timestamp (8 bytes) | modified timestamp (8 bytes) | data length (4 bytes) | data (variable) |
func writeEntry(w io.Writer, key string, checksum uint32, created, modified int64, data []byte) error {
	// Write key length and key
	keyLen := uint32(len(key))
	if err := binary.Write(w, binary.LittleEndian, keyLen); err != nil {
		return err
	}
	if _, err := w.Write([]byte(key)); err != nil {
		return err
	}

	// Write checksum
	if err := binary.Write(w, binary.LittleEndian, checksum); err != nil {
		return err
	}

	// Write created timestamp
	if err := binary.Write(w, binary.LittleEndian, created); err != nil {
		return err
	}

	// Write modified timestamp
	if err := binary.Write(w, binary.LittleEndian, modified); err != nil {
		return err
	}

	// Write data length and data
	dataLen := uint32(len(data))
	if err := binary.Write(w, binary.LittleEndian, dataLen); err != nil {
		return err
	}
	_, err := w.Write(data)

	return err
}

// readEntry reads a key-value pair with length prefixes
// Same format as writeEntry
// Returns the key, data, and any error encountered
func readEntry(r io.Reader, skipData bool) (*ObjectMeta, []byte, error) {
	// Read key length and key
	var keyLen uint32
	if err := binary.Read(r, binary.LittleEndian, &keyLen); err != nil {
		return nil, nil, err
	}

	keyBuf := make([]byte, keyLen)
	if _, err := io.ReadFull(r, keyBuf); err != nil {
		return nil, nil, err
	}
	key := string(keyBuf)

	// Read checksum
	var checksum uint32
	if err := binary.Read(r, binary.LittleEndian, &checksum); err != nil {
		return nil, nil, err
	}

	// Read created timestamp
	var created int64
	if err := binary.Read(r, binary.LittleEndian, &created); err != nil {
		return nil, nil, err
	}

	// Read modified timestamp
	var modified int64
	if err := binary.Read(r, binary.LittleEndian, &modified); err != nil {
		return nil, nil, err
	}

	// Read data length and data
	var dataLen uint32
	if err := binary.Read(r, binary.LittleEndian, &dataLen); err != nil {
		return nil, nil, err
	}

	if skipData {
		// Skip the data by seeking forward
		if seeker, ok := r.(io.Seeker); ok {
			if _, err := seeker.Seek(int64(dataLen), io.SeekCurrent); err != nil {
				return nil, nil, err
			}
			return &ObjectMeta{
				Key:        key,
				Offset:     0, // Offset is not relevant when skipping data
				Size:       dataLen,
				Checksum:   checksum,
				CreatedAt:  created,
				ModifiedAt: modified,
			}, nil, nil
		} else {
			// If the reader doesn't support seeking, read and discard the data
			if _, err := io.CopyN(io.Discard, r, int64(dataLen)); err != nil {
				return nil, nil, err
			}
			return &ObjectMeta{
				Key:        key,
				Offset:     0, // Offset is not relevant when skipping data
				Size:       dataLen,
				Checksum:   checksum,
				CreatedAt:  created,
				ModifiedAt: modified,
			}, nil, nil
		}
	}

	dataBuf := make([]byte, dataLen)
	if _, err := io.ReadFull(r, dataBuf); err != nil {
		return nil, nil, err
	}

	meta := &ObjectMeta{
		Key:        key,
		Offset:     0,
		Size:       dataLen,
		Checksum:   checksum,
		CreatedAt:  created,
		ModifiedAt: modified,
	}

	return meta, dataBuf, nil
}

func readAllEntries(s *shard) ([]*ObjectMeta, error) {
	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	currentOffset, err := s.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	var entries []*ObjectMeta
	for {
		meta, _, err := readEntry(s.file, true)
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		meta.Offset = currentOffset

		entries = append(entries, meta)

		currentOffset, err = s.file.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
	}

	return entries, nil
}
