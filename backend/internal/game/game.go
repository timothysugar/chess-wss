package game

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/notnil/chess"
	"github.com/timothysugar/chess/pkg/websocket"
)

type Player struct {
	Colour   chess.Color
	NextTurn bool
	client   *websocket.Client
}

func NewPlayer(client *websocket.Client, colour chess.Color) *Player {
	return &Player{
		client: client,
		Colour: colour,
	}
}

type Match struct {
	Id         string
	pool       []*websocket.Client
	Players    map[chess.Color]*Player
	Game       *chess.Game
	NextColour chess.Color

	Unregister chan *Player
	Moves      chan string

	Broadcast chan string

	ticker *time.Ticker
}

func NewMatch() *Match {
	return &Match{
		Id:         "001",
		Players:    make(map[chess.Color]*Player, 2),
		NextColour: chess.White,
	}
}

func (m *Match) broadcastValidMoves() {
	body, err := json.Marshal(m.Game.ValidMoves())
	if err != nil {
		fmt.Println("unable to encode game moves to JSON %v\n", err)
	}
	m.broadcastMessage(string(body))
}

func (m *Match) begin() {
	m.Game = chess.NewGame()
	m.NextColour = chess.White

	m.Players[chess.White] = NewPlayer(m.pool[0], chess.White)
	m.Players[chess.Black] = NewPlayer(m.pool[1], chess.Black)
	m.broadcastMessage("Starting match")
	whiteMoves := m.Players[chess.White].client.Read()
	whiteOut := func(msg string) error {
		return m.Players[chess.White].client.Write(msg)
	}
	blackMoves := m.Players[chess.Black].client.Read()
	blackOut := func(msg string) error {
		return m.Players[chess.Black].client.Write(msg)
	}

	for {
		select {
		case msg, ok := <-whiteMoves:
			fmt.Println("received from white", msg, ok)
			err := m.Game.MoveStr(msg.Body)
			if err != nil {
				fmt.Println("invalid move received from white %v\n", err)
				whiteOut("invalid move")
			} else {
				fmt.Println("white moved")
			}
			m.broadcastValidMoves()
		case msg, ok := <-blackMoves:
			fmt.Println("received from black", msg, ok)
			err := m.Game.MoveStr(msg.Body)
			if err != nil {
				fmt.Printf("invalid move received from black %v\n", err)
				blackOut("invalid move")
			} else {
				fmt.Println("black moved")
			}
			m.broadcastValidMoves()
		default:
			fmt.Println("waiting for move")
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func (m *Match) Join(client *websocket.Client) error {
	fmt.Printf("client joined the game %v %v\n", client, m.pool)
	num := len(m.pool)
	if num >= 2 {
		return fmt.Errorf("%d players already in match %s", num, m.Id)
	}

	m.broadcastMessage("Player joining")
	m.pool = append(m.pool, client)

	if len(m.pool) == 2 {
		m.begin()
	}

	return nil
}

func (m *Match) Move(move string) error {
	err := m.Game.MoveStr(move)
	m.broadcastMessage(fmt.Sprintf("Move made %s", move))
	return err
}

func (m *Match) broadcastMessage(msg string) {
	for _, client := range m.pool {
		err := client.Write(msg)
		if err != nil {
			fmt.Printf("Error sending broadcast message to player %v %v", client, err)
		}
	}
}

// func (m *Match) tick() {
// 	m.broadcastMessage(time.Now().Format(time.RFC1123))
// }

// func (m *Match) play() {
// 	for {
// 		select {
// 		case _ = <-m.ticker.C:
// 			m.tick()
// 		case client := <-m.Unregister:
// 			delete(m.Players, client)
// 		case msg := <-m.Broadcast:
// 			m.broadcastMessage(msg)
// 		}
// 	}
// }
