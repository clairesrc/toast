package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	addr string
	cors string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// allow localhost:4000
		if r.Header.Get("Origin") == "http://localhost:4000" {
			return true
		}

		// allow all origins
		return true
	},
} // use default options

func (wss WebsocketServer) echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", wss.cors)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (wss WebsocketServer) start() error {
	http.HandleFunc("/echo", wss.echo)
	err := http.ListenAndServe(wss.addr, nil)
	if err != nil {
		return fmt.Errorf("ListenAndServe: %v", err)
	}

	return nil
}
