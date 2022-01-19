package main

import (
	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/go-common/server"
	"github.com/webdevelop-pro/notification-worker/internal/adapters/mapstorage"
	"github.com/webdevelop-pro/notification-worker/internal/app"
	"github.com/webdevelop-pro/notification-worker/internal/ports/http"
	"go.uber.org/fx"
)

// @schemes https
func main() {
	fx.New(
		fx.Logger(logger.NewDefaultComponent("fx")),
		fx.Provide(
			// Default logger
			logger.NewDefault,
			// Repository
			mapstorage.New,
			// Init http server
			server.New,
			// Init application
			app.New,
		),
		// Init http handlers
		http.InitServer(),

		fx.Invoke(
			// Run HTTP server
			server.StartServer,
		),
	).Run()
}
