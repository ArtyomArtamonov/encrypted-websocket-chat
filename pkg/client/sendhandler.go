package client

import (
	"bufio"
	"os"
)

func (c *Client) handleSender(messageCh chan string) {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := scanner.Text()
		messageCh <- line
	}
}
