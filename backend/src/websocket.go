package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	addr      string
	cors      string
	gameState *gameState
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
}

func (wss WebsocketServer) state(w http.ResponseWriter, r *http.Request) {
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

		// parse the message from json to gameEvent
		event := gameEvent{}
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Println("json unmarshal:", err)
			break
		}
		wss.gameState.handleEvent(event)

		err = c.WriteMessage(mt, wss.gameState.toJSON())
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (wss WebsocketServer) start() error {
	http.HandleFunc("/state", wss.state)
	fmt.Println("Websocket server starting on", wss.addr)
	err := http.ListenAndServe(wss.addr, nil)
	if err != nil {
		return fmt.Errorf("ListenAndServe: %v", err)
	}

	fmt.Println("Websocket server started on", wss.addr)

	return nil
}
