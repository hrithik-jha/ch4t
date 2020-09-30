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

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func handleMessages() {
	for {
		// Acquiring next message
		msg := <-Broadcast

		// Sending to every client
		for Client := range Clients {
			err := Client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error: %v", err)
				Client.Close()
				delete(Clients, Client)
			}
		}
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	Clients[ws] = true
	for {
		var msg Message

		// Reading message as JSON and mapping to Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error: %v", err)
			delete(Clients, ws)
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
