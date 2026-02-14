package client

import (
	"encoding/binary"

	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/protocol"
)

// Ping implements the client side of the protocol.ServerCommandPing command.
func (c *Client) Ping() error {
	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandPing,
		DatabaseName:   "",
		CollectionName: "",
		Payload:        make([]byte, 0),
	})
	if err != nil {
		return err
	}

	if resp.Code != protocol.StatusOK {
		return ErrUnexpectedStatusCode
	}

	return nil
}

// GetSupportedProtocolVersions implements the client side of the protocol.ServerCommandSupportedProtocolVersions command.
func (c *Client) GetSupportedProtocolVersions() ([]byte, error) {
	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandSupportedProtocolVersions,
		DatabaseName:   "",
		CollectionName: "",
		Payload:        make([]byte, 0),
	})
	if err != nil {
		return nil, err
	}

	if resp.Code != protocol.StatusOK {
		return nil, ErrUnexpectedStatusCode
	}

	numVersions := resp.Payload[0]
	versions := make([]byte, numVersions)
	copy(versions, resp.Payload[1:1+numVersions])

	return versions, nil
}

// LoginWithPassword implements the client side of the protocol.ServerCommandLogin command for password-based authentication.
func (c *Client) LoginWithPassword(username, password string) error {
	totalLen := 1 + 2 + len(username) + 2 + len(password)
	payload := make([]byte, totalLen)

	payload[0] = 0x01 // type: password

	binary.BigEndian.PutUint16(payload[1:3], uint16(len(username))) // username length (2 bytes)
	copy(payload[3:3+len(username)], []byte(username))              // username

	binary.BigEndian.PutUint16(payload[3+len(username):5+len(username)], uint16(len(password))) // password length (2 bytes)
	copy(payload[5+len(username):], []byte(password))                                           // password

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandLogin,
		DatabaseName:   "",
		CollectionName: "",
		Payload:        payload,
	})
	if err != nil {
		return err
	}

	if resp.Code != protocol.StatusOK {
		if resp.Code == protocol.StatusInvalidCredentials {
			return ErrInvalidCredentials
		}

		return ErrUnexpectedStatusCode
	}

	return nil
}

// LoginWithApiKey implements the client side of the protocol.ServerCommandLogin command for API key-based authentication.
func (c *Client) LoginWithApiKey(apiKey string) error {
	totalLen := 1 + 2 + len(apiKey)
	payload := make([]byte, totalLen)

	payload[0] = 0x02 // type: apiKey

	binary.BigEndian.PutUint16(payload[1:3], uint16(len(apiKey))) // apiKey length (2 bytes)
	copy(payload[3:], []byte(apiKey))                             // apiKey

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandLogin,
		DatabaseName:   "",
		CollectionName: "",
		Payload:        payload,
	})
	if err != nil {
		return err
	}

	if resp.Code != protocol.StatusOK {
		if resp.Code == protocol.StatusInvalidCredentials {
			return ErrInvalidCredentials
		}

		return ErrUnexpectedStatusCode
	}

	return nil
}

// IsAuthenticated implements the client side of the protocol.ServerCommandAuthStatus command.
func (c *Client) IsAuthenticated() (bool, error) {
	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandAuthStatus,
		DatabaseName:   "",
		CollectionName: "",
		Payload:        make([]byte, 0),
	})
	if err != nil {
		return false, err
	}

	if resp.Code == protocol.StatusOK {
		return true, nil
	}

	if resp.Code == protocol.StatusUnauthorized {
		return false, nil
	}

	return false, ErrUnexpectedStatusCode
}
