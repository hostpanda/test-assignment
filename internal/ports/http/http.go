package http

import (
	httpHandler "github.com/webdevelop-pro/notification-worker/internal/ports/addHTTP"
	wsHandler "github.com/webdevelop-pro/notification-worker/internal/ports/addWS"
	"go.uber.org/fx"
)

// InitServer returns a new instance of Server
func InitServer() fx.Option {
	return fx.Options(
		fx.Invoke(
			// Registration routes and handlers for http server
			wsHandler.InitHandler,
			httpHandler.InitHandler,
		),
	)
}
