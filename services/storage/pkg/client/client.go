package client

import (
	"log/slog"
	"net"
	"strconv"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

type Client struct {
	conn net.Conn
}

type Configuration struct {
	Host string
	Port string

	Username, Password string
	ApiKey             string
}

func NewClient(cfg Configuration) (*Client, error) {
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == "" {
		cfg.Port = "8091"
	}

	conn, err := net.Dial("tcp", cfg.Host+":"+cfg.Port)
	if err != nil {
		return nil, err
	}

	c := &Client{conn: conn}

	if err := c.Ping(); err != nil {
		c.Close()
		return nil, err
	}

	if err := c.checkProtocolVersionSupport(); err != nil {
		return nil, err
	}

	if cfg.ApiKey != "" {
		if err := c.LoginWithApiKey(cfg.ApiKey); err != nil {
			c.Close()
			return nil, err
		}
	} else if cfg.Username != "" && cfg.Password != "" {
		if err := c.LoginWithPassword(cfg.Username, cfg.Password); err != nil {
			c.Close()
			return nil, err
		}
	}

	slog.Info("FancySpaces storage client connected", slog.String("host", cfg.Host), slog.String("port", cfg.Port))
	return c, nil
}

func (c *Client) Close() {
	c.conn.Close()
	c.conn = nil
	return
}

func (c *Client) IsConnected() bool {
	return c.conn != nil
}

func (c *Client) checkProtocolVersionSupport() error {
	supportedProtocolVersions, err := c.GetSupportedProtocolVersions()
	if err != nil {
		c.Close()
		return err
	}

	supportsCurrentVersion := false
	for _, ver := range supportedProtocolVersions {
		if ver == byte(protocol.V1.Version) {
			supportsCurrentVersion = true
		}
	}

	if !supportsCurrentVersion {
		c.Close()
		return ErrProtocolVersionNotSupported
	}

	return nil
}

func (c *Client) SendCmd(cmd *protocol.Command) (*protocol.Response, error) {
	if !c.IsConnected() {
		return nil, ErrClientNotConnected
	}

	msg := protocol.Message{
		ProtocolVersion: byte(protocol.ProtocolVersion1),
		Flags:           0x00,
		Type:            byte(protocol.MessageTypeCommand),
		Payload:         protocol.V1.EncodeCommand(cmd),
	}

	data := protocol.V1.EncodeMessage(&msg)
	if err := protocol.V1.WriteFrame(c.conn, data); err != nil {
		return nil, err
	}

	slog.Debug(
		"Sent command to server",
		slog.String("id", strconv.Itoa(int(cmd.ID))),
		slog.String("database", cmd.DatabaseName),
		slog.String("collection", cmd.CollectionName),
		slog.String("payload_size", strconv.Itoa(len(cmd.Payload))),
	)

	return c.readResponse()
}

func (c *Client) readResponse() (*protocol.Response, error) {
	if !c.IsConnected() {
		return nil, ErrClientNotConnected
	}

	frame, err := protocol.V1.ReadFrame(c.conn)
	if err != nil {
		return nil, err
	}
	msg, err := protocol.V1.DecodeMessage(frame)
	if err != nil {
		return nil, err
	}
	resp, err := protocol.V1.DecodeResponse(msg)
	if err != nil {
		return nil, err
	}

	slog.Debug(
		"Received response from server",
		slog.String("status_code", strconv.Itoa(int(resp.Code))),
		slog.String("payload_size", strconv.Itoa(len(resp.Payload))),
	)

	return resp, nil
}
