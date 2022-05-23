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

		if c.aesKey == nil {
			c.marshalAndSendRecursivly(MessageFrameType, []byte(line), c.sendRSA)
			continue
		}
		
		c.marshalAndSendRecursivly(MessageFrameType, []byte(line), c.sendAES)
	}
}

func (c *Client) sendAES(data []byte) {
	if c.aesKey != nil {
		encryptedPayload, err := ciphers.EncryptDataAES(data, c.aesKey)
		if err != nil {
			c.sendAES(data[:len(data) / 2])
			c.sendAES(data[len(data)/2:])
		}
		c.conn.WriteMessage(1, encryptedPayload)
		return
	}

	c.conn.WriteMessage(1, data)
}

func (c *Client) marshalAndSendRecursivly(frameType FrameType, data []byte, recursiveSendFunc func([]byte)) {
	frame := NewFrame(c.name, data, frameType)

	payload, err := json.Marshal(frame)
	if err != nil {
		log.Fatalf("Could not marshal frame %v", frame)
	}

	recursiveSendFunc(payload)
}

func (c *Client) sendRSA(data []byte) {
	if c.partnerKeys.PublicKey != nil {
		encryptedPayload, err := ciphers.EncryptWithPublicKey(data, c.partnerKeys.PublicKey)
		if err != nil {
			c.sendRSA(data[:len(data) / 2])
			c.sendRSA(data[len(data)/2:])
		}
		c.conn.WriteMessage(1, encryptedPayload)
		return
	}

	c.conn.WriteMessage(1, data)
}
