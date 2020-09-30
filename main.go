package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	clients[ws] = true
	for {
		var msg Message

		// Reading message as JSON and mapping to Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error: %v", err)
			delete(clients, ws)
			break
		}

		// Send the message to Message Queue
		Broadcast <- msg
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/home", Home)
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	log.Println("Starting server on PORT:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
