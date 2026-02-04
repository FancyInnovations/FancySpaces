package command

import (
	"encoding/binary"

	"github.com/fancyinnovations/fancyspaces/storage/internal/auth"
	"github.com/fancyinnovations/fancyspaces/storage/internal/protocol"
)

var (
	pongResponse = &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: []byte("pong"),
	}

	supportedProtocolVersionsResponse = &protocol.Response{
		Code: protocol.StatusOK,
		Payload: []byte{
			1, // length of supported versions
			byte(protocol.ProtocolVersion1),
		},
	}

	invalidCredentialsResponse = &protocol.Response{
		Code:    protocol.StatusInvalidCredentials,
		Payload: *emptyPayload,
	}

	emptyOK = &protocol.Response{
		Code:    protocol.StatusOK,
		Payload: *emptyPayload,
	}

	unsupportedAuthTypeResponse = &protocol.Response{
		Code:    protocol.StatusBadRequest,
		Payload: []byte("unsupported authentication type"),
	}

	unauthorizedResponse = &protocol.Response{
		Code:    protocol.StatusUnauthorized,
		Payload: *emptyPayload,
	}
)

var emptyPayload = &[]byte{}

// handlePing processes a ping command.
func handlePing(_ *ConnCtx, _ *protocol.Message, _ *protocol.Command) (*protocol.Response, error) {
	return pongResponse, nil
}

// handleSupportedProtocolVersions processes a supported protocol versions command.
func handleSupportedProtocolVersions(_ *ConnCtx, _ *protocol.Message, _ *protocol.Command) (*protocol.Response, error) {
	return supportedProtocolVersionsResponse, nil
}

// handleLogin processes a login command.
func handleLogin(ctx *ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
	data := cmd.Payload
	if len(data) == 0 {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("authentication type is required"),
		}, nil
	}

	authType := data[0]
	switch authType {
	case 1: // Password
		return handleLoginWithPassword(ctx, nil, cmd)
	case 2: // API key
		return handleLoginWithAPIKey(ctx, nil, cmd)
	default:
		return unsupportedAuthTypeResponse, nil
	}
}

func handleLoginWithPassword(ctx *ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
	data := cmd.Payload

	if len(data) < 3 {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload for basic authentication"),
		}, nil
	}

	usernameLen := binary.BigEndian.Uint16(data[1:3])
	if len(data) < int(3+usernameLen+2) {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload for basic authentication"),
		}, nil
	}
	username := string(data[3 : 3+usernameLen])

	passwordLenStart := 3 + usernameLen
	passwordLen := binary.BigEndian.Uint16(data[passwordLenStart : passwordLenStart+2])
	if len(data) < int(passwordLenStart+2+passwordLen) {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload for basic authentication"),
		}, nil
	}
	password := string(data[passwordLenStart+2 : passwordLenStart+2+passwordLen])

	newCtx, err := auth.AuthenticateWithBasicAuth(ctx.Ctx, username, password)
	if err != nil {
		return invalidCredentialsResponse, nil
	}
	ctx.Ctx = newCtx

	return emptyOK, nil
}

func handleLoginWithAPIKey(ctx *ConnCtx, _ *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
	data := cmd.Payload

	if len(data) < 3 {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload for API key authentication"),
		}, nil
	}

	apiKeyLen := binary.BigEndian.Uint16(data[1:3])
	if len(data) < int(3+apiKeyLen) {
		return &protocol.Response{
			Code:    protocol.StatusBadRequest,
			Payload: []byte("invalid payload for API key authentication"),
		}, nil
	}
	apiKey := string(data[3 : 3+apiKeyLen])

	newCtx, err := auth.AuthenticateWithApiKey(ctx.Ctx, apiKey)
	if err != nil {
		return invalidCredentialsResponse, nil
	}
	ctx.Ctx = newCtx

	return emptyOK, nil
}

// handleAuthStatus processes an auth status command.
func handleAuthStatus(ctx *ConnCtx, _ *protocol.Message, _ *protocol.Command) (*protocol.Response, error) {
	user := auth.UserFromContext(ctx.Ctx)
	if user == nil {
		return unauthorizedResponse, nil
	}

	return emptyOK, nil
}
