package client

import (
	"log/slog"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

type Client struct {
	cfg              Configuration
	conn             net.Conn
	requestIDCounter uint32
	pendingCmds      map[uint32]chan *protocol.Response
	pendingCmdsMu    sync.Mutex
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
		cfg:              cfg,
		conn:             conn,
		pendingCmds:      make(map[uint32]chan *protocol.Response),
		requestIDCounter: 0,
	}

	go c.startResponseListener()

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

	c.pendingCmdsMu.Lock()
	c.pendingCmds = make(map[uint32]chan *protocol.Response)
	c.pendingCmdsMu.Unlock()

	go c.startResponseListener()

	if err := c.Ping(); err != nil {
		slog.Error("Ping failed after reconnecting", slog.Any("error", err))
		c.Close()
		return
	}

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

	respChan := make(chan *protocol.Response, 1)
	c.pendingCmdsMu.Lock()
	cmd.ReqID = c.getNextRequestID()
	c.pendingCmds[cmd.ReqID] = respChan
	c.pendingCmdsMu.Unlock()

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

	// Wait for the response or a timeout
	select {
	case resp := <-respChan:
		duration := time.Since(starTime)
		slog.Debug(
			"Received response from server",
			slog.String("status_code", strconv.Itoa(int(resp.Code))),
			slog.String("payload_size", strconv.Itoa(len(resp.Payload))),
			slog.Duration("duration", duration),
		)
		return resp, nil
	case <-time.After(30 * time.Second):
		c.pendingCmdsMu.Lock()
		delete(c.pendingCmds, cmd.ReqID)
		c.pendingCmdsMu.Unlock()
		return nil, ErrCommandTimeout
	}
}
