package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()
	var addr = flag.String("addr", ":8181", "http service address")

	wss := WebsocketServer{
		addr: *addr,
		cors: "*",
		gameState: &gameState{
			Players: []player{},
		},
	}
	err := wss.start()
	if err != nil {
		log.Fatal(err)
	}
}
