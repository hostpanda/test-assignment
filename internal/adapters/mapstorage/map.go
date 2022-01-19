package mapstorage

import (
	"context"
	"encoding/json"
	"fmt"

	comLogger "github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/notification-worker/internal/adapters"
)

// DB is a layer to simplify interact with DB
type MapStorage struct {
	data map[string]json.RawMessage
	log  comLogger.Logger
}

// New returns new DB instance.
func New() adapters.Repository {
	return NewMapStorage(comLogger.NewDefaultComponent("mapstorage"))
}

// NewMapStorage returns new DB instance.
func NewMapStorage(log comLogger.Logger) *MapStorage {
	ms := &MapStorage{
		data: map[string]json.RawMessage{},
		log:  log,
	}

	return ms
}

// Save element in the map
func (ms *MapStorage) Save(ctx context.Context, number int) error {
	return fmt.Errorf("implement save")
}
