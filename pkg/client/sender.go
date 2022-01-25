package client

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/ArtyomArtamonov/encrypted-websocket-chat/pkg/ciphers"
)

func (c *Client) sender() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := scanner.Text()

		c.marshalAndSend(MessageFrameType, []byte(line))
	}
}

func (c *Client) marshalAndSend(_type FrameType, data []byte) {
	frame := NewFrame(c.name, data, _type)

	payload, err := json.Marshal(frame)
	if err != nil {
		log.Fatalf("Could not marshal frame %v", frame)
	}

	c.send(payload)
}

func (c *Client) send(data []byte) {
	if c.partnerKeys.PublicKey != nil {
		encryptedPayload, err := ciphers.EncryptWithPublicKey(data, c.partnerKeys.PublicKey)
		if err != nil {
			c.send(data[:len(data) / 2])
			c.send(data[len(data)/2:])
		}
	
		c.conn.WriteMessage(1, encryptedPayload)
		return
	}

	c.conn.WriteMessage(1, data)
}
