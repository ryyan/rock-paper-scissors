package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/grafov/bcast"
	"golang.org/x/net/websocket"
)

var (
	// Command line args
	port = flag.String("port", ":5000", "Service port")

	// Global broadcast channel
	// Used to trigger websocket pushes
	broadcaster = bcast.NewGroup()

	currentGame = &Game{
		Left:  0,
		Right: 0,
	}
	currentState = &GameState{
		LeftTaken:  false,
		RightTaken: false,
		Wins:       []int64{0, 0, 0},
		Ties:       []int64{0, 0, 0},
	}
)

type Game struct {
	Left  int64 // 0=None, 1=Rock, 10=Paper, 100=Scissors
	Right int64
	sync.Mutex
}

type GameState struct {
	LeftTaken  bool
	RightTaken bool
	Wins       []int64 // Rock, Paper, Scissor wins
	Ties       []int64 // Rock, Paper, Scissor ties
	sync.Mutex
}

func rpsHandler(w http.ResponseWriter, r *http.Request) {
	qp := r.URL.Query()
	leftOrRight := qp.Get("lor") // Left=l, Right=r
	choice := qp.Get("choice")   // Rock=1, Paper=10, Scissors=100

	// Validate left or right
	if leftOrRight != "l" && leftOrRight != "r" {
		w.WriteHeader(401)
		w.Write([]byte("Invalide left/right"))
	}

	// Validate choice
	choiceInt, err := strconv.ParseInt(choice, 10, 64)
	if err != nil ||
		(choiceInt != 1 && choiceInt != 10 && choiceInt != 100) {
		w.WriteHeader(400)
		w.Write([]byte("Invalid choice"))
	}

	// Check if left or right is already taken
	currentState.Lock()
	defer currentState.Unlock()
	if (leftOrRight == "l" && currentState.LeftTaken) ||
		(leftOrRight == "r" && currentState.RightTaken) {
		w.WriteHeader(401)
		w.Write([]byte("Left/Right already taken"))
	}

	// Lock in choice
	currentGame.Lock()
	defer currentGame.Unlock()
	if leftOrRight == "l" {
		currentGame.Left = int64(choiceInt)
		currentState.LeftTaken = true
	} else {
		currentGame.Right = int64(choiceInt)
		currentState.RightTaken = true
	}

	// Perform game logic
	gameWasPlayed := true
	switch currentGame.Left + currentGame.Right {
	case 2:
		// Rock vs Rock
		currentState.Ties[0] = currentState.Ties[0] + 1
	case 20:
		// Paper vs Paper
		currentState.Ties[1] = currentState.Ties[1] + 1
	case 200:
		// Scissors vs Scissors
		currentState.Ties[2] = currentState.Ties[2] + 1
	case 11:
		// Rock vs Paper
		currentState.Wins[1] = currentState.Wins[1] + 1
	case 101:
		// Rock vs Scissors
		currentState.Wins[0] = currentState.Wins[0] + 1
	case 110:
		// Paper vs Scissors
		currentState.Wins[2] = currentState.Wins[2] + 1
	default:
		gameWasPlayed = false
	}

	// Reset current game if one was played
	if gameWasPlayed {
		currentGame = &Game{
			Left:  0,
			Right: 0,
		}
	}

	// Trigger broadcast for websocket listeners
	broadcaster.Send(true)
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func rpsWebsocket(ws *websocket.Conn) {
	// Send initial response
	message, _ := json.Marshal(currentState)
	io.WriteString(ws, string(message))

	// Send responses on change
	for {
		listener := broadcaster.Join()
		listener.Recv()
		message, _ := json.Marshal(currentState)
		io.WriteString(ws, string(message))
	}
}

func main() {
	// Parse command line args
	flag.Parse()

	// Start broadcaster
	go broadcaster.Broadcast(2 * time.Minute)

	// Start server
	http.HandleFunc("/rps", rpsHandler)
	http.Handle("/websocket/rps", websocket.Handler(rpsWebsocket))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("../client/public"))))
	log.Println("Serving at localhost" + *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}
