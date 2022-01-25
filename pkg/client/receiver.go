package client

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ArtyomArtamonov/encrypted-websocket-chat/pkg/ciphers"
)

func (c *Client) receiver(doneCh chan struct{}) {
	defer close(doneCh)
	publicKeyWasSent := false
	messageConstructor := ""

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Fatal("read:", err)
			return
		}

		decrypted, err := ciphers.DecryptWithPrivateKey(message, c.PrivateKey)
		if err != nil {
			decrypted = message
		}
		
		frame := Frame{}
		err = json.Unmarshal(decrypted, &frame)
		if err != nil {
			messageConstructor += string(decrypted)
			errMsg := json.Unmarshal([]byte(messageConstructor), &frame)
			if errMsg != nil {
				continue
			}
			messageConstructor = ""
		}

		switch frame.Type{
		case MessageFrameType:
			fmt.Printf("%s: %s\n", frame.From, string(frame.Data))
		case SystemFrameType:
			log.Printf("Got system message: %v", frame.Data)
		case HandshakeFrameType:
			key := ciphers.BytesToPublicKey(frame.Data)
			c.partnerKeys.PublicKey = key
			
			if !publicKeyWasSent {
				c.marshalAndSend(HandshakeFrameType, ciphers.PublicKeyToBytes(c.PublicKey))
				publicKeyWasSent = true
			}
		}
	}
}
