package client

import (
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

type Client struct {
	cfg  Configuration
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

	c := &Client{
		cfg:  cfg,
		conn: conn,
	}

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

	go c.startHeartbeat()

	slog.Info("FancySpaces storage client connected", slog.String("host", cfg.Host), slog.String("port", cfg.Port))
	return c, nil
}

func (c *Client) Reconnect() {
	c.Close()

	conn, err := net.Dial("tcp", c.cfg.Host+":"+c.cfg.Port)
	if err != nil {
		slog.Error("Failed to reconnect to server", slog.String("host", c.cfg.Host), slog.String("port", c.cfg.Port), slog.Any("error", err))
		return
	}

	c.conn = conn

	slog.Info("FancySpaces storage client reconnected", slog.String("host", c.cfg.Host), slog.String("port", c.cfg.Port))
}

func (c *Client) Close() {
	if c.conn == nil {
		return
	}

	c.conn.Close()
	c.conn = nil

	slog.Info("FancySpaces storage client disconnected", slog.String("host", c.cfg.Host), slog.String("port", c.cfg.Port))
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

func (c *Client) startHeartbeat() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := c.Ping(); err != nil {
			slog.Warn("Heartbeat ping failed, reconnecting", slog.Any("error", err))
			c.Reconnect()
			continue
		}

		slog.Debug("Heartbeat ping successful")
	}
}

func (c *Client) SendCmd(cmd *protocol.Command) (*protocol.Response, error) {
	if !c.IsConnected() {
		return nil, ErrClientNotConnected
	}

	starTime := time.Now()

	cmdMsg := protocol.Message{
		ProtocolVersion: byte(protocol.ProtocolVersion1),
		Flags:           0x00,
		Type:            byte(protocol.MessageTypeCommand),
		Payload:         protocol.V1.EncodeCommand(cmd),
	}

	cmdData := protocol.V1.EncodeMessage(&cmdMsg)
	if err := protocol.V1.WriteFrame(c.conn, cmdData); err != nil {
		return nil, err
	}

	slog.Debug(
		"Sent command to server",
		slog.String("id", strconv.Itoa(int(cmd.ID))),
		slog.String("database", cmd.DatabaseName),
		slog.String("collection", cmd.CollectionName),
		slog.String("payload_size", strconv.Itoa(len(cmd.Payload))),
	)

	// Read response
	respFrame, err := protocol.V1.ReadFrame(c.conn)
	if err != nil {
		return nil, err
	}
	respMsg, err := protocol.V1.DecodeMessage(respFrame)
	if err != nil {
		return nil, err
	}
	resp, err := protocol.V1.DecodeResponse(respMsg)
	if err != nil {
		return nil, err
	}

	duration := time.Since(starTime)
	slog.Debug(
		"Received response from server",
		slog.String("status_code", strconv.Itoa(int(resp.Code))),
		slog.String("payload_size", strconv.Itoa(len(resp.Payload))),
		slog.Duration("duration", duration),
	)

	return resp, nil
}
