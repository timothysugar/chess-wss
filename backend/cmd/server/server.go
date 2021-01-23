package main

import (
	"fmt"
	"net/http"

	"github.com/timothysugar/chess/internal/game"
	"github.com/timothysugar/chess/pkg/websocket"
)

func serveWs(match *game.Match, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := websocket.Client{
		Conn: conn,
	}
	match.Join(&client)
}

func setupRoutes() {
	match := game.NewMatch()
	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		serveWs(match, w, r)
	})
}

func main() {
	fmt.Println("Distributed Chat App v0.01")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
