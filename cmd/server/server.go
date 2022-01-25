package main

import (
	"flag"

	"github.com/ArtyomArtamonov/encrypted-websocket-chat/pkg/server"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "http service address")

	server := server.NewServer(1024, 1024)
	server.Run(*addr)
}

