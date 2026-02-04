package command

import "github.com/fancyinnovations/fancyspaces/storage/internal/protocol"

var pongResponse = &protocol.Response{
	Code:    protocol.StatusOK,
	Payload: []byte("pong"),
}

func handlePing(_ *protocol.Message, _ *protocol.Command) (*protocol.Response, error) {
	return pongResponse, nil
}

var supportedProtocolVersionsResponse = &protocol.Response{
	Code: protocol.StatusOK,
	Payload: []byte{
		1, // length of supported versions
		byte(protocol.ProtocolVersion1),
	},
}

func handleSupportedProtocolVersions(_ *protocol.Message, _ *protocol.Command) (*protocol.Response, error) {
	return supportedProtocolVersionsResponse, nil
}
