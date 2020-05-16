package pingo_pongo

import (
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

const htmlTemplate = "game.html"

var wsUpgrader = websocket.Upgrader{}
var gamePool = NewGamePool()

type PlayerServer struct {
	http.Handler

	template *template.Template
}

func NewPlayerServer() (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles(htmlTemplate)
	if err != nil {
		return nil, fmt.Errorf("problem opening %s, %w", htmlTemplate, err)
	}

	p.template = tmpl

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(p.index))
	router.Handle("/ws", http.HandlerFunc(p.webSocketHandler))

	p.Handler = router

	return p, nil
}

func (g *PlayerServer) index(w http.ResponseWriter, r *http.Request) {
	if err := g.template.Execute(w, nil); err != nil {
		http.Error(w, fmt.Sprintf("problem executing template %v", err), http.StatusInternalServerError)
	}
}

func (g *PlayerServer) webSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := newServerWS(w, r)
	if err != nil {
		log.Printf("problem create WebSocket server, %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player, err := NewPlayer(conn)
	if err != nil {
		log.Printf("problem create player, %v\n", err)
		return
	}

	if err := gamePool.EnterTheGame(player); err != nil {
		log.Printf("problem enter the game, %v\n", err)
		return
	}
}
