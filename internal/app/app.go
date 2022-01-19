package app

import (
	"context"

	"github.com/webdevelop-pro/go-common/configurator"
	"github.com/webdevelop-pro/go-common/logger"

	"github.com/webdevelop-pro/notification-worker/internal/adapters"
)

type Application struct {
	config     *Config
	log        logger.Logger
	sessionHub *sessionHub
	repo       adapters.Repository
}

func New(repo adapters.Repository) *Application {
	cfg := &Config{}
	log := logger.NewDefaultComponent("app")

	if err := configurator.NewConfiguration(cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to get configuration of app")
	}

	a := &Application{
		config:     cfg,
		log:        log,
		sessionHub: newSessionHub(),
		repo:       repo,
	}

	return a
}

func (a *Application) NewConnection(ctx context.Context, userID string, conn adapters.Conn) error {

	NewSession(ctx, userID, conn, a.sessionHub)

	return nil
}
