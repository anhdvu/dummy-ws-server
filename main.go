package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type wsserver struct {
	upgrader websocket.Upgrader
	clients  map[string]*websocket.Conn
}

func NewWSServer() *wsserver {
	fn := func(r *http.Request) bool {
		return true
	}

	return &wsserver{
		upgrader: websocket.Upgrader{
			CheckOrigin: fn,
		},
		clients: make(map[string]*websocket.Conn),
	}
}

func (s *wsserver) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", s.homeHandler())
	mux.Handle("/ws", s.wsEndpointHandler())

	return mux
}

func (s *wsserver) homeHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "aulabs home")
	})
}

func (s *wsserver) wsEndpointHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		clientID := uuid.New()
		s.clients[clientID.String()] = ws
		err = ws.WriteMessage(1, []byte("Welcome "+clientID.String()))
	})
}

func (s *wsserver) spamMessage() {
	messages := []string{
		"hello trym",
		"hello any",
		"hello aulabs",
		"to the moon",
		"luna ftw",
		"sc2 coop, go?",
		"apm 300",
	}

	const interval = 2

	for {
		for clientID, clientConn := range s.clients {
			m := fmt.Sprintf("client ID: %s\n%s\n", clientID, random(messages...))
			err := clientConn.WriteMessage(1, []byte(m))
			if err != nil {
				log.Printf("error %s sending message to client %s", err.Error(), clientID)
				delete(s.clients, clientID)
			}
			time.Sleep(interval * time.Second)
		}
	}

}

func random(msg ...string) string {
	l := len(msg)

	if l == 0 {
		return ""
	}

	return msg[rand.Intn(l)]
}

func main() {
	srv := NewWSServer()

	server := &http.Server{
		Addr:         ":7878",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      srv.routes(),
	}

	go srv.spamMessage()

	log.Fatal(server.ListenAndServe())
}
