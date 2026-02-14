package dbcmds

import (
	"errors"
	"log/slog"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/commonresponses"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/protocol"
	"github.com/fancyinnovations/fancyspaces/storage/internal/auth"
	"github.com/fancyinnovations/fancyspaces/storage/internal/command"
	"github.com/fancyinnovations/fancyspaces/storage/internal/database"
)

// PermissionMiddleware is a middleware that checks if the user has the necessary permissions to execute a command on a database.
func (c *Commands) PermissionMiddleware(next command.Handler) command.Handler {
	return func(ctx *command.ConnCtx, msg *protocol.Message, cmd *protocol.Command) (*protocol.Response, error) {
		if cmd.DatabaseName == "" {
			return next(ctx, msg, cmd)
		}

		u := auth.UserFromContext(ctx.Ctx)
		if u == nil || !u.Verified || !u.IsActive {
			return commonresponses.Unauthorized, nil
		}

		if u.IsAdmin() {
			return next(ctx, msg, cmd) // Admin users have access to all databases, so we can skip the permission check for them.
		}

		db, err := c.dbStore.GetDatabase(ctx.Ctx, cmd.DatabaseName)
		if err != nil {
			if errors.Is(err, database.ErrDatabaseNotFound) {
				return commonresponses.DatabaseNotFound, nil
			}

			slog.Error("Failed to get database",
				slog.String("database", cmd.DatabaseName),
				sloki.WrapError(err),
			)
			return commonresponses.InternalServerError, nil
		}

		if !db.HasPermission(u.ID, database.PermissionLevelReadWrite) {
			return commonresponses.Forbidden, nil
		}

		return next(ctx, msg, cmd)
	}
}
