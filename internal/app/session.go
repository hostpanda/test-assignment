package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/webdevelop-pro/notification-worker/internal/domain/ws"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/notification-worker/internal/adapters"
)

const (
	sessionWriteWait      = 10 * time.Second
	sessionPongWait       = 60 * time.Second
	sessionPingPeriod     = (sessionPongWait * 9) / 10
	sessionMaxMessageSize = 10000
)

var (
	ErrSessionClosed = errors.New("tried to write to a closed session")

	pingMessage  = &message{t: websocket.PingMessage, data: nil}
	closeMessage = &message{t: websocket.CloseMessage, data: nil}
	null         void
)

//type sessionConfig struct {
//	WriteWait         time.Duration `default:"10s" split_words:"true"`   // Time allowed to write a message to the peer.
//	PongWait          time.Duration `default:"1m" split_words:"true"`    // Time allowed to read the next pong message from the peer.
//	PingPeriod        time.Duration `default:"55s" split_words:"true"`   // Send pings to peer with this period. Must be less than pongWait.
//	MaxMessageSize    int64         `default:"10000" split_words:"true"` // Maximum message size allowed from client.
//	MessageBufferSize int           `default:"256" split_words:"true"`   // The max amount of messages that can be in a sessions buffer before it starts dropping them.
//}

type void struct{}

type message struct {
	t    int
	data []byte
}

type Payload struct {
	Event   string          `json:"event,omitempty"`
	Channel string          `json:"channel,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type signalToHub interface {
	RemoveSession(s *session)
	AddSession(s *session)
}

type session struct {
	log      logger.Logger
	conn     adapters.Conn
	userID   string
	channels map[string]void
	output   chan *message
	closed   bool
	rwMu     *sync.RWMutex
	hub      signalToHub
}

func NewSession(ctx context.Context, userID string, conn adapters.Conn, hub signalToHub) *session {
	s := &session{
		log:      logger.Logger{Logger: *log.Ctx(ctx)},
		conn:     conn,
		userID:   userID,
		channels: make(map[string]void),
		output:   make(chan *message),
		closed:   false,
		rwMu:     &sync.RWMutex{},
		hub:      hub,
	}
	s.log.Trace().Msg("new session")

	s.hub.AddSession(s)
	go s.readPump()
	go s.writePump()

	return s
}

func (s *session) isClosed() bool {
	s.rwMu.RLock()
	defer s.rwMu.RUnlock()
	return s.closed
}

func (s *session) write(msg *message) error {
	if s.isClosed() {
		return ErrSessionClosed
	}

	_ = s.conn.SetWriteDeadline(time.Now().Add(sessionWriteWait))

	if err := s.conn.WriteMessage(msg.t, msg.data); err != nil {
		return err
	}

	return nil
}

func (s *session) ping() error {
	return s.write(pingMessage)
}

// Close closes session.
func (s *session) Close() error {
	return s.write(closeMessage)
}

func (s *session) writePump() {
	ticker := time.NewTicker(sessionPingPeriod)

	defer func() {
		ticker.Stop()
		s.close()
	}()

	for {
		select {
		case msg, ok := <-s.output:
			if !ok {
				// The hub already closed the channel.
				return
			}

			if msg.t == websocket.CloseMessage {
				return
			}

			if err := s.write(msg); err != nil {
				return
			}

		case <-ticker.C:
			if err := s.ping(); err != nil {
				return
			}
		}
	}
}

func (s *session) readPump() {
	//defer c.Close()
	s.conn.SetReadLimit(sessionMaxMessageSize)
	_ = s.conn.SetReadDeadline(time.Now().Add(sessionPongWait))

	s.conn.SetPongHandler(
		func(string) error {
			return s.conn.SetReadDeadline(time.Now().Add(sessionPongWait))
		},
	)

	for {
		t, msg, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.log.Error().Err(err).Msg("failed to read an incoming message")
			}
			break
		}

		switch t {
		case websocket.TextMessage:
			s.handleTextMessage(msg)
		default:
			s.log.Error().Err(err).Msg(fmt.Sprintf("unexpected websocket type: %d", t))
		}

	}
}

func (s *session) handleTextMessage(payload []byte) {
	s.log.Trace().Msg(fmt.Sprintf("%s", payload))
	var pld ws.Payload

	if err := json.Unmarshal(payload, &pld); err != nil {
		log.Error().Err(err).Msg("failed to decode message")
		return
	}

	switch pld.Event {
	case ws.SubscribeEvent:
		// {"event":"subscribe","channels":["test"]}
		s.rwMu.Lock()
		for i := range pld.Channels {
			s.channels[pld.Channels[i]] = null
		}
		s.rwMu.Unlock()
	case ws.UnsubscribeEvent:
		s.rwMu.Lock()
		for i := range pld.Channels {
			delete(s.channels, pld.Channels[i])
		}
		s.rwMu.Unlock()
	}

	return
}

func (s *session) close() {
	s.rwMu.Lock()
	defer s.rwMu.Unlock()

	if s.closed {
		return
	}

	s.hub.RemoveSession(s)
	close(s.output)

	if err := s.conn.Close(); err != nil {
		s.log.Error().Err(err).Msg("failed to close ws connection")
	}

	s.closed = true

}
