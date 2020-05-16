package pingo_pongo

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	writeTimeout = 1 * time.Second
	readTimeout  = 1 * time.Second
)

type serverWS struct {
	*websocket.Conn
	mutex sync.Mutex
}

func newServerWS(w http.ResponseWriter, r *http.Request) (*serverWS, error) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("problem upgrading connection to WebSockets, %v", err)
		return nil, err
	}

	conn.SetCloseHandler(func(code int, text string) error {
		message := websocket.FormatCloseMessage(code, text)
		return conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(writeTimeout))
	})

	return &serverWS{Conn: conn}, nil
}

func (s *serverWS) send(v interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
		return fmt.Errorf("send failed, %w", err)
	}

	if err := s.WriteJSON(v); err != nil {
		return fmt.Errorf("send failed, %w", err)
	}

	return nil
}

func (s *serverWS) receive(v interface{}) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.SetReadDeadline(time.Now().Add(readTimeout)); err != nil {
		return fmt.Errorf("receive failed, %w", err)
	}

	if err := s.ReadJSON(v); err != nil {
		return fmt.Errorf("receive failed, %w", err)
	}

	return nil
}

func (s *serverWS) ping() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeTimeout)); err != nil {
		return fmt.Errorf("ping failed, %w", err)
	}

	return nil
}
