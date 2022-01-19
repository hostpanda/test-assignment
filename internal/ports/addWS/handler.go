package wsHandler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/go-common/server"
	"github.com/webdevelop-pro/notification-worker/internal/app"
)

// Handler represents http handler for investments
type Handler struct {
	app      *app.Application
	log      logger.Logger
	upgrader websocket.Upgrader
}

// InitHandler initiates a new investment handler
func InitHandler(srv *server.HttpServer, app *app.Application) {
	h := &Handler{
		app: app,
		log: logger.NewDefaultComponent("ws_handler"),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			// TODO get origin host from env
			//CheckOrigin: func(r *http.Request) bool {
			//	origin := r.Header["Origin"]
			//	if len(origin) == 0 {
			//		return false
			//	}
			//	u, err := url.Parse(origin[0])
			//	if err != nil {
			//		return false
			//	}
			//	return strings.HasSuffix(u.Host, "")
			//},
		},
	}

	for _, route := range h.GetRoutes() {
		srv.AddRoute(route)
	}
}

func (h *Handler) GetRoutes() []server.Route {

	return []server.Route{
		{
			Method: http.MethodGet,
			Path:   "/addWS/",
			Handle: h.WSHandler,
			NoAuth: true,
		},
	}
}
