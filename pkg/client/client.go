package client

import (
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/ArtyomArtamonov/encrypted-websocket-chat/pkg/ciphers"
	"github.com/gorilla/websocket"
)

type Client struct {
	url url.URL
	conn *websocket.Conn
	name string

	ciphers.CryptoKeys
	partnerKeys ciphers.CryptoKeys
}

func NewClient(url url.URL, name string) *Client {
	return &Client{
		url: url,
		name: name,
	}
}

const KEY_LENGTH = 2048

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

	// Generate and send key before messaging
	{
		priv, pub := ciphers.GenerateKeyPair(KEY_LENGTH)
		cl.PrivateKey = priv
		cl.PublicKey = pub

		cl.marshalAndSend(RsaHandshakeFrameType, ciphers.PublicKeyToBytes(pub))
	}
	
	
	// handleReceiver will handle incoming messages
	done := make(chan struct{})
	go cl.receiver(done)

	// handleSender will handle outgoing messages 
	// provided by stdin
	go cl.sender()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupt")

			cl.marshalAndSend(PartnerDisconnectedFrameType, []byte{})

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
