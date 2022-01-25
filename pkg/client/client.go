package client

import (
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

type Client struct {
	url url.URL
	conn *websocket.Conn
	name string
}

func NewClient(url url.URL, name string) *Client {
	return &Client{
		url: url,
		name: name,
	}
}

func (cl Client) Run() {
	// Ctrl+C notifier for graceful stopping
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("Connecting to %s", cl.url.String())

	c, _, err := websocket.DefaultDialer.Dial(cl.url.String(), nil)
	if err != nil {
		log.Fatal("Dial:", err)
	}
	defer c.Close()
	cl.conn = c

	// handleReceiver will handle incoming messages
	done := make(chan struct{})
	go cl.handleReceive(done)

	// handleSender will handle outgoing messages 
	// provided by stdin
	var messageChan = make(chan string)
	go cl.handleSender(messageChan)

	for {
		select {
		case msg := <-messageChan:
			c.WriteMessage(1, []byte(msg))
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Write close:", err)
				return
			}
			return
		}
	}
}
