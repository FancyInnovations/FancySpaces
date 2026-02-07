package objectengine

import (
	"encoding/binary"
	"io"
)

// writeEntry writes a key-value pair with length prefixes
func writeEntry(w io.Writer, key string, data []byte) error {
	// Write key length and key
	keyLen := uint32(len(key))
	if err := binary.Write(w, binary.LittleEndian, keyLen); err != nil {
		return err
	}
	if _, err := w.Write([]byte(key)); err != nil {
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
func readEntry(r io.Reader) (key string, data []byte, err error) {
	// Read key length and key
	var keyLen uint32
	if err = binary.Read(r, binary.LittleEndian, &keyLen); err != nil {
		return "", nil, err
	}

	keyBuf := make([]byte, keyLen)
	if _, err = io.ReadFull(r, keyBuf); err != nil {
		return "", nil, err
	}
	key = string(keyBuf)

	// Read data length and data
	var dataLen uint32
	if err = binary.Read(r, binary.LittleEndian, &dataLen); err != nil {
		return "", nil, err
	}

	data = make([]byte, dataLen)
	if _, err = io.ReadFull(r, data); err != nil {
		return "", nil, err
	}

	return key, data, nil
}
