package client

import (
	"encoding/json"
	"fmt"
	"log"
)

func (c *Client) handleReceive(doneCh chan struct{}) {
	defer close(doneCh)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Fatal("read:", err)
			return
		}
		
		frame := Frame{}
		err = json.Unmarshal(message, &frame)
		if err != nil {
			log.Fatal("Could not unmarshall message", message)
		}

		switch frame.Type{
		case SystemFrameType:
			log.Printf("Got system message: %v", frame.Data)
		case MessageFrameType:
			fmt.Printf("%s: %s\n", frame.From, fmt.Sprintf("%s", frame.Data))
		}
	}
}
