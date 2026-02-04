package protocol

const magicNumber byte = 0x7E

type Message struct {
	ProtocolVersion byte
	Flags           byte
	Type            byte
	Payload         []byte
}

type MessageType byte

const (
	MessageTypeRequest  MessageType = 0x01
	MessageTypeResponse MessageType = 0x02
	MessageTypeError    MessageType = 0x03
)

type ProtocolVersion byte

const (
	ProtocolVersion1 ProtocolVersion = 0x01
)

// TODO: implement flags as bitmask
//type MessageFlag byte
//
//const (
//	MessageFlagCompressed MessageFlag = 1 << 0
//	MessageFlagEncrypted  MessageFlag = 1 << 1
//)

type Command struct {
	ID             uint16
	DatabaseName   string
	CollectionName string
	Payload        []byte
}
