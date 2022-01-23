package main

import "github.com/ArtyomArtamonov/encrypted-websocket-chat/pkg/server"

func main() {
	server := server.NewServer(1024, 1024)
	server.Run("localhost:8080")
}

