package protocol

import (
	"sync"
)

var (
	requestBufferPool = sync.Pool{
		New: func() any {
			return make([]byte, 0, 1024) // TODO: tune initial size based on typical request sizes
		},
	}

	responseBufferPool = sync.Pool{
		New: func() any {
			return make([]byte, 0, 1024) // TODO: tune initial size based on typical request sizes
		},
	}

	messagePool = sync.Pool{
		New: func() any {
			return &Message{}
		},
	}

	commandPool = sync.Pool{
		New: func() any {
			return &Command{}
		},
	}
)

// GetRequestBufferFromPool retrieves a byte slice from the request buffer pool.
// Use this for encoding request payloads.
func GetRequestBufferFromPool() []byte {
	return requestBufferPool.Get().([]byte)[:0]
}

func PutRequestBufferToPool(buf []byte) {
	// avoid holding on to huge allocations
	if cap(buf) > 64*1024 {
		return
	}
	requestBufferPool.Put(buf[:0])
}

// GetResponseBufferFromPool retrieves a byte slice from the response buffer pool.
// Use this for encoding response payloads.
func GetResponseBufferFromPool() []byte {
	return responseBufferPool.Get().([]byte)[:0]
}

func PutResponseBufferToPool(buf []byte) {
	// avoid holding on to huge allocations
	if cap(buf) > 64*1024 { // TODO: tune this threshold based on typical response sizes
		return
	}
	responseBufferPool.Put(buf[:0])
}

func GetMessageFromPool() *Message {
	return messagePool.Get().(*Message)
}

func PutMessageToPool(m *Message) {
	// reset fields
	m.ProtocolVersion = 0
	m.Flags = 0
	m.Type = 0
	m.Payload = nil

	messagePool.Put(m)
}

func GetCommandFromPool() *Command {
	return commandPool.Get().(*Command)
}

func PutCommandToPool(c *Command) {
	// reset fields
	c.ID = 0
	c.DatabaseName = ""
	c.CollectionName = ""
	c.Payload = nil

	commandPool.Put(c)
}
