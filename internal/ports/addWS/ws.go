package wsHandler

import (
	"github.com/labstack/echo/v4"
	"github.com/webdevelop-pro/go-common/server/middleware"
)

func (h *Handler) WSHandler(c echo.Context) error {
	conn, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to upgrade connection")
		return err
	}

	ctx := c.Request().Context()
	jwtPayload := middleware.GetJWTPayload(ctx)

	return h.app.NewConnection(ctx, jwtPayload.UserID, conn)
}
