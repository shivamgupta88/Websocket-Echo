package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var messageHistory []string
var historyMux sync.Mutex

func reverseMessage(msg string) string {
	runes := []rune(msg)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}

		message := string(msg)
		fmt.Println("Received:", message)

		reversed := reverseMessage(message)

		historyMux.Lock()
		messageHistory = append(messageHistory, message)
		if len(messageHistory) > 5 {
			messageHistory = messageHistory[1:]
		}
		historyMux.Unlock()

		err = conn.WriteMessage(websocket.TextMessage, []byte(reversed))
		if err != nil {
			log.Println("Failed to send message:", err)
			break
		}
	}
}
