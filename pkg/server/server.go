package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader *websocket.Upgrader
	storage *ConnectionStorage
}

func NewServer(readBufferSize, writeBufferSize int) *Server {
	return &Server{
		upgrader: &websocket.Upgrader{
			ReadBufferSize: readBufferSize,
			WriteBufferSize: writeBufferSize,
		},
		storage: NewStorage(),
	}
}

func (s *Server) Run(ipaddr string) {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Print("Got new connection")
		conn, err := s.upgrader.Upgrade(rw, r, nil)
		if err != nil {
			fmt.Fprint(rw, "Could not upgrade")
			log.Fatal("Could not upgrade")
		}
		s.storage.add(conn)
		
		conn.SetCloseHandler(s.getCloseHandler(conn))
		go s.handleMessages(conn)
		
	})

	log.Print("Server started")
	http.ListenAndServe(ipaddr, nil)
}

func (s *Server) getCloseHandler(conn *websocket.Conn) func(code int, text string) error {
	handler := func(_ int, _ string) error {
		s.storage.delete(conn)
		return nil
	}

	return handler
}

func (s *Server) handleMessages(c *websocket.Conn){
	log.Print("Server listening for client messages")
	for {
		// Read message from clients
		msgType, msg, err := c.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", c.RemoteAddr(), string(msg))

		s.storage.Lock()
		for user := range s.storage.Connections {
			if user == c {
				continue
			}
			if err = user.WriteMessage(msgType, msg); err != nil {
				log.Printf("Could not send message to user")
				continue
			}
		}
		s.storage.Unlock()
	}
}
