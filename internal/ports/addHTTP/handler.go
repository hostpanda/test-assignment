package httpHandler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/go-common/server"
	"github.com/webdevelop-pro/notification-worker/internal/app"
	"github.com/webdevelop-pro/notification-worker/internal/domain/event"
)

// Handler represents http handler for investments
type Handler struct {
	app *app.Application
	log logger.Logger
}

// InitHandler initiates a new investment handler
func InitHandler(srv *server.HttpServer, app *app.Application) {
	h := &Handler{
		app: app,
		log: logger.NewDefaultComponent("http_handler"),
	}

	for _, route := range h.GetRoutes() {
		srv.AddRoute(route)
	}
}

func (h *Handler) CreateRecord(ctx echo.Context) error {
	e := event.Event{}
	if err := ctx.Bind(&e); err != nil {
		h.log.Error().Err(err).Interface("ev", e).Msg("cannot process record")
		return ctx.String(400, "cannot process request")
	}
	fmt.Println("we need to save record")
	return nil
}

func (h *Handler) GetRoutes() []server.Route {

	return []server.Route{
		{
			Method: http.MethodPost,
			Path:   "/addHTTP/",
			Handle: h.CreateRecord,
			NoAuth: true,
		},
	}
}
