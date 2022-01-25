package client

import (
	"fmt"
	"log"
)

func (c *Client) handleReceive(doneCh chan struct{}) {
	defer close(doneCh)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		fmt.Printf("Message: %s\n", message)
	}
}
