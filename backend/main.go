package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	log.Println("loaded words:", len(wordList))
	http.HandleFunc("/ws", handleWs)
	fmt.Println("Go Websockets")

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
