package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	go func() {
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error on reading the message %v", err)
				return
			}

			fmt.Printf("WS Message: %s", p)

			if err := conn.WriteMessage(messageType, p); err != nil {
				log.Printf("Error on responding to message %v", err)
				return
			}
		}
	}()
	go func() {
		for {
			err := conn.WriteMessage(websocket.TextMessage, []byte("^Ping\n"))
			if err != nil {
				log.Printf("Error on sending the message %v", err)
			}
			fmt.Print("Message sent")
			time.Sleep(5 * time.Second)
		}
	}()
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error on upgrading the socket %v", err)
		return
	}

	fmt.Println("Websocket connected")
	reader(ws)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().Second())
	fmt.Println(time.Now().Second()%5 == 0)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}

func setupRoutes() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
