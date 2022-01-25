package main

import (
	"flag"
	"net/url"

	"github.com/ArtyomArtamonov/encrypted-websocket-chat/pkg/client"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "http service address")
	name := flag.String("name", "username", "name that will be seen")
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	client := client.NewClient(u, *name)

	client.Run()
}
