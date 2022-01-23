package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/ArtyomArtamonov/encrypted-websocket-chat/pkg/client"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	client := client.NewClient(u)

	client.Run()
}
