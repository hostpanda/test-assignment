package adapters

import (
	"io"
	"time"
)

type Conn interface {
	SetPongHandler(h func(appData string) error)

	SetReadLimit(limit int64)
	SetReadDeadline(t time.Time) error

	SetWriteDeadline(t time.Time) error

	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)

	NextWriter(messageType int) (io.WriteCloser, error)

	Close() error
}
