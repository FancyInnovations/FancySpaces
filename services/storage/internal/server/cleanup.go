package server

import (
	"log/slog"
	"time"
)

const DisconnectAfterInactivity = 60 * 1000 // 1 minute

func (s *Server) cleanupInactiveConnections() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		slog.Info("Running cleanup of inactive connections")
		now := time.Now().UnixMilli()

		s.connectionsMu.Lock()
		for id, connCtx := range s.connections {
			if now-connCtx.LastActivity > DisconnectAfterInactivity {
				slog.Info("Disconnecting inactive connection", slog.String("connection_id", id))
				connCtx.Conn.Close()
				delete(s.connections, id)
			}
		}
		s.connectionsMu.Unlock()
	}
}
