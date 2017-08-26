package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/grafov/bcast"
	"golang.org/x/net/websocket"
)

var (
	// Command line args
	port = flag.String("port", ":5000", "Service port")

	// Global broadcast channel
	// Used to trigger websocket pushes
	broadcaster = bcast.NewGroup()

	// Game state
	currentGame = &Game{
		Left:          0,
		Right:         0,
		LeftTaken:     false,
		RightTaken:    false,
		Wins:          []int64{0, 0, 0},
		Ties:          []int64{0, 0, 0},
		PreviousGames: []*GameRecord{},
	}

	// Used to marhsall game state to json once and
	// write to all websockets
	currentGameJson, _ = json.Marshal(currentGame)
)

type Game struct {
	Left          int64 `json:"-"` // 0=None, 1=Rock, 10=Paper, 100=Scissors
	Right         int64 `json:"-"` // "-" to prevent returning as json
	LeftTaken     bool
	RightTaken    bool
	Wins          []int64 // Rock, Paper, Scissor wins
	Ties          []int64 // Rock, Paper, Scissor ties
	PreviousGames []*GameRecord
	sync.RWMutex
}

type GameRecord struct {
	Left  string
	Right string
}

func rpsHandler(w http.ResponseWriter, r *http.Request) {
	// Lock current game and game state
	fmt.Println("Lock")
	currentGame.Lock()
	defer currentGame.Unlock()
	defer fmt.Println("Unlock")

	// Parse query para
	qp := r.URL.Query()
	leftOrRight := qp.Get("lor") // Left=l, Right=r
	choice := qp.Get("choice")   // Rock=1, Paper=10, Scissors=100

	// Validate left or right
	if leftOrRight != "l" && leftOrRight != "r" {
		w.WriteHeader(401)
		w.Write([]byte("Invalide left/right"))
		return
	}

	// Validate choice
	choiceInt, err := strconv.ParseInt(choice, 10, 64)
	if err != nil ||
		(choiceInt != 1 && choiceInt != 10 && choiceInt != 100) {
		w.WriteHeader(400)
		w.Write([]byte("Invalid choice"))
		return
	}

	// Check if left or right is already taken
	if (leftOrRight == "l" && currentGame.LeftTaken) ||
		(leftOrRight == "r" && currentGame.RightTaken) {
		w.WriteHeader(401)
		w.Write([]byte("Left/Right already taken"))
		return
	}

	// Lock in choice
	if leftOrRight == "l" {
		currentGame.Left = int64(choiceInt)
		currentGame.LeftTaken = true
	} else {
		currentGame.Right = int64(choiceInt)
		currentGame.RightTaken = true
	}

	// Perform game logic
	gameWasPlayed := true
	switch currentGame.Left + currentGame.Right {
	case 2:
		// Rock vs Rock
		currentGame.Ties[0] = currentGame.Ties[0] + 1
	case 20:
		// Paper vs Paper
		currentGame.Ties[1] = currentGame.Ties[1] + 1
	case 200:
		// Scissors vs Scissors
		currentGame.Ties[2] = currentGame.Ties[2] + 1
	case 11:
		// Rock vs Paper
		currentGame.Wins[1] = currentGame.Wins[1] + 1
	case 101:
		// Rock vs Scissors
		currentGame.Wins[0] = currentGame.Wins[0] + 1
	case 110:
		// Paper vs Scissors
		currentGame.Wins[2] = currentGame.Wins[2] + 1
	default:
		gameWasPlayed = false
	}

	// Record game and reset if one was played
	if gameWasPlayed {
		record := &GameRecord{
			Left:  choiceToString(currentGame.Left),
			Right: choiceToString(currentGame.Right),
		}
		currentGame.PreviousGames = append(currentGame.PreviousGames, record)

		// Only keep the last 10 games
		if len(currentGame.PreviousGames) == 11 {
			currentGame.PreviousGames = currentGame.PreviousGames[1:len(currentGame.PreviousGames)]
		}

		currentGame.Left = 0
		currentGame.Right = 0
		currentGame.LeftTaken = false
		currentGame.RightTaken = false
	}

	// Trigger broadcast for websocket listeners
	currentGameJson, _ = json.Marshal(currentGame)
	defer broadcaster.Send(true)
	w.WriteHeader(200)
	w.Write([]byte("OK"))
	return
}

func rpsWebsocket(ws *websocket.Conn) {
	// Send initial response
	io.WriteString(ws, string(currentGameJson))

	// Send responses on change
	for {
		listener := broadcaster.Join()
		listener.Recv()
		io.WriteString(ws, string(currentGameJson))
	}
}

func choiceToString(choice int64) string {
	switch choice {
	case 1:
		return "Rock"
	case 10:
		return "Paper"
	case 100:
		return "Scissors"
	default:
		return ""
	}
}

func main() {
	// Parse command line args
	flag.Parse()

	// Start broadcaster
	go broadcaster.Broadcast(0)

	// Start server
	http.HandleFunc("/rps", rpsHandler)
	http.Handle("/websocket/rps", websocket.Handler(rpsWebsocket))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("../client/public"))))
	log.Println("Serving at localhost" + *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}
