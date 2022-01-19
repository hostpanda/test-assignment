package app

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/webdevelop-pro/go-common/logger"
	"github.com/webdevelop-pro/notification-worker/internal/domain/event"
)

type sessionHub struct {
	sessions map[string]map[*session]bool
	rwMu     *sync.RWMutex
	log      logger.Logger
}

func newSessionHub() *sessionHub {
	return &sessionHub{
		sessions: make(map[string]map[*session]bool),
		rwMu:     &sync.RWMutex{},
		log:      logger.NewDefaultComponent("session_hub"),
	}
}

func (sh *sessionHub) send(msg event.Event) {
	sh.rwMu.RLock()
	defer sh.rwMu.RUnlock()

	log := sh.log.With().Interface("message", msg).Logger()

	switch msg.UserID {
	case "":
		payload, err := json.Marshal(msg)
		if err != nil {
			log.Trace().Msg("failed marshal message")

			return
		}

		for userID := range sh.sessions {
			for s := range sh.sessions[userID] {
				if _, ok := s.channels[msg.Channel]; !ok {
					continue
				}

				ses := s
				go func() {
					ses.output <- &message{t: websocket.TextMessage, data: payload}
				}()
			}
		}

	default:
		sessions, ok := sh.sessions[msg.UserID]
		if !ok {
			log.Trace().Msg("session not found")

			return
		}

		for s := range sessions {

			if _, ok := s.channels[msg.Channel]; !ok {
				return
			}

			ses := s
			go func() {
				payload, err := json.Marshal(msg)
				if err != nil {
					log.Trace().Msg("failed marshal message")

					return
				}

				ses.output <- &message{t: websocket.TextMessage, data: payload}
			}()

		}
	}
}

func (sh *sessionHub) RemoveSession(s *session) {
	s.log.Trace().Msg("remove session")
	go func() {
		sh.rwMu.Lock()
		defer sh.rwMu.Unlock()

		sessions, ok := sh.sessions[s.userID]
		if !ok {
			return
		}

		countSessions := len(sessions)

		delete(sessions, s)

		if countSessions == 1 {
			delete(sh.sessions, s.userID)
		}

	}()
}

func (sh *sessionHub) AddSession(s *session) {
	s.log.Trace().Msg("add session")
	go func() {
		sh.rwMu.Lock()
		defer sh.rwMu.Unlock()

		sessions, ok := sh.sessions[s.userID]
		if !ok {
			sessions = make(map[*session]bool)
			sh.sessions[s.userID] = sessions
		}

		sessions[s] = true
	}()
}

// Send sends message to the subscriber. If userID is empty broadcasts to all subscribers.
func (sh *sessionHub) Send(msg event.Event) {
	go sh.send(msg)
}

func (sh *sessionHub) Close() {
	sh.log.Trace().Msg("Close session hub")
	go func() {
		sh.rwMu.Lock()
		defer sh.rwMu.Unlock()

		for userID := range sh.sessions {
			for s := range sh.sessions[userID] {
				_ = s.Close()

				countSessions := len(sh.sessions[userID])

				delete(sh.sessions[userID], s)

				if countSessions == 1 {
					delete(sh.sessions, s.userID)
				}
			}
		}
	}()
}
