package client

import (
	"errors"
	"log/slog"
	"net"
	"strconv"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

func (c *Client) getNextRequestID() uint32 {
	c.requestIDCounter++

	if c.requestIDCounter >= 0xFFFFFFFF {
		c.requestIDCounter = 1
	}

	return c.requestIDCounter
}

func (c *Client) startResponseListener() {
	for {
		if !c.IsConnected() {
			slog.Debug("Connection closed, stopping response listener")
			return
		}

		respFrame, err := protocol.V1.ReadFrame(c.conn)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				slog.Debug("Connection closed, stopping response listener")
				return
			}

			slog.Warn("Failed to read response frame", slog.Any("error", err))
			continue
		}

		respMsg, err := protocol.V1.DecodeMessage(respFrame)
		if err != nil {
			slog.Warn("Failed to decode response message", slog.Any("error", err))
			continue
		}

		if respMsg.Type == byte(protocol.MessageTypeResponse) {
			resp, err := protocol.V1.DecodeResponse(respMsg)
			if err != nil {
				slog.Warn("Failed to decode response", slog.Any("error", err))
				continue
			}

			c.pendingCmdsMu.Lock()
			respChan, exists := c.pendingCmds[resp.ReqID]
			if exists {
				delete(c.pendingCmds, resp.ReqID)
				c.pendingCmdsMu.Unlock()
				respChan <- resp
			} else {
				c.pendingCmdsMu.Unlock()
				slog.Warn("Received response with unknown ID", slog.String("response_id", strconv.Itoa(int(resp.ReqID))))
			}
		} else if respMsg.Type == byte(protocol.MessageTypeCommand) {
			cmd, err := protocol.V1.DecodeCommand(respMsg)
			if err != nil {
				slog.Warn("Failed to decode command message", slog.Any("error", err))
				continue
			}

			slog.Info("Received command from server",
				slog.String("command_id", strconv.Itoa(int(cmd.ID))),
				slog.String("database", cmd.DatabaseName),
				slog.String("collection", cmd.CollectionName),
				slog.String("payload_size", strconv.Itoa(len(cmd.Payload))),
			)
		} else {
			slog.Warn("Received message with unknown type", slog.String("message_type", strconv.Itoa(int(respMsg.Type))))
		}
	}
}
