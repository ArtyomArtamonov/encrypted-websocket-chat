package client

import (
	"bufio"
	"os"
)

func (c *Client) handleSender(frameCh chan *Frame) {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := scanner.Text()

		frame := NewFrame(c.name, line, MessageFrameType)
		frameCh <- frame
	}
}
