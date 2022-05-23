package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionStorage struct {
	Connections map[*websocket.Conn]bool
	sync.Mutex
}

func NewStorage() *ConnectionStorage {
	return &ConnectionStorage{
		Connections: make(map[*websocket.Conn]bool),
	}
}

func (cs *ConnectionStorage) add(conn *websocket.Conn) {
	cs.Lock()
	defer cs.Unlock()

	cs.Connections[conn] = true
}

func (cs *ConnectionStorage) delete(conn *websocket.Conn) {
	cs.Lock()
	defer cs.Unlock()

	delete(cs.Connections, conn)
}
