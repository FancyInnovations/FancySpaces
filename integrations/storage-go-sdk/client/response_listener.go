package client

import (
	"encoding/binary"
	"errors"
	"log/slog"
	"net"
	"strconv"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
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

			// Handle broker messages
			if cmd.ID == protocol.ClientCommandBrokerMessage {
				payload := cmd.Payload
				if len(payload) < 2+8 { // subject length (2 bytes) + at least 8 bytes for list
					slog.Warn("Received broker message with invalid payload", slog.String("payload_size", strconv.Itoa(len(payload))))
					continue
				}

				subjectLen := int(binary.BigEndian.Uint16(payload[0:2]))
				if len(payload) < 2+subjectLen+8 {
					slog.Warn(
						"Received broker message with invalid payload (subject length mismatch)",
						slog.String("payload_size", strconv.Itoa(len(payload))),
						slog.String("subject_length", strconv.Itoa(subjectLen)),
					)
					continue
				}
				subject := string(payload[2 : 2+subjectLen])

				msgs, err := codex.DecodeList(payload[2+subjectLen:])
				if err != nil {
					slog.Warn("Failed to decode broker message list", slog.Any("error", err))
					continue
				}

				c.brokerSubjectListenersMu.Lock()
				listeners := c.brokerSubjectListeners[cmd.DatabaseName+"."+cmd.CollectionName+"."+subject]
				c.brokerSubjectListenersMu.Unlock()

				for _, msg := range msgs {
					for _, listener := range listeners {
						go listener(msg.AsBinary())
					}
				}
			}
		} else {
			slog.Warn("Received message with unknown type", slog.String("message_type", strconv.Itoa(int(respMsg.Type))))
		}
	}
}
