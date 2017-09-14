package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/grafov/bcast"
	"golang.org/x/net/websocket"
)

var (
	// Command line args
	port = flag.Int("port", 5000, "Service port")

	// Global broadcast channels
	// Used to trigger websocket pushes on game state changes
	gameBroadcaster    = bcast.NewGroup()
	playersBroadcaster = bcast.NewGroup()

	currentPlayers = int64(0)

	currentGame = &Game{
		Left:          0,
		Right:         0,
		LeftTaken:     false,
		RightTaken:    false,
		Wins:          []int64{0, 0, 0},
		Ties:          []int64{0, 0, 0},
		PreviousGames: []*GameRecord{},
	}
)

// Game holds the global game state
type Game struct {
	// Left/Right hold the int value of rock/paper/scissor
	Left  int64 `json:"-"` // "-" to prevent returning as json
	Right int64 `json:"-"`

	// LeftTaken/RightTaken are true if Left/Right are non-zero
	// This is returned to the client instead of Left/Right to keep choices hidden
	LeftTaken  bool
	RightTaken bool

	// [Rock, Paper, Scissors] wins/ties
	Wins []int64
	Ties []int64

	// Previous games
	PreviousGames []*GameRecord

	// Mutex to ensure only one client changes game state at a time
	sync.RWMutex
}

type GameRecord struct {
	Left  string
	Right string
}

// rpsHandler handles game input and logic
func rpsHandler(w http.ResponseWriter, r *http.Request) {
	// Lock current game and game state
	currentGame.Lock()
	defer currentGame.Unlock()

	// Parse query parameters
	qp := r.URL.Query()
	leftOrRight := qp.Get("lor") // Left=l, Right=r
	choice := qp.Get("choice")   // Rock=1, Paper=10, Scissors=100

	// Validate left or right
	if leftOrRight != "l" && leftOrRight != "r" {
		w.WriteHeader(401)
		w.Write([]byte("Invalid left/right"))
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

	// Return response for this call
	w.WriteHeader(200)
	w.Write([]byte("OK"))

	// Do not trigger websocket broadcast until this call is returned
	// otherwise a deadlock will occur from this client
	gameBroadcaster.Send(true)
}

// rpsWebsocket handler sends Game state json to clients
func rpsWebsocketHandler(ws *websocket.Conn) {
	// Send initial response
	err := websocket.JSON.Send(ws, currentGame)
	if err != nil {
		return
	}

	// Send responses on change
	listener := gameBroadcaster.Join()
	for {
		listener.Recv() // Blocks until gameBroadcaster receives a message
		err := websocket.JSON.Send(ws, currentGame)
		if err != nil {
			return
		}
	}
}

// currentPlayersWebsocket sends current number of active players (ws connections)
func currentPlayersWebsocketHandler(ws *websocket.Conn) {
	// Increment current players on ws connection
	atomic.AddInt64(&currentPlayers, 1)
	playersBroadcaster.Send(true)

	// Send initial response
	err := websocket.JSON.Send(ws, currentPlayers)
	if err != nil {
		atomic.AddInt64(&currentPlayers, -1)
		playersBroadcaster.Send(true)
		return
	}

	go func() {
		listener := playersBroadcaster.Join()
		for {
			listener.Recv()
			err := websocket.JSON.Send(ws, currentPlayers)
			if err != nil {
				return
			}
		}
	}()

	// Decrement current players on ws disconnect
	websocket.Message.Receive(ws, nil)
	atomic.AddInt64(&currentPlayers, -1)
	playersBroadcaster.Send(true)
}

// choiceToString converts an int to its matching choice string
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

func newWebsocketHandler(handler websocket.Handler) websocket.Handler {
	// Override the websocket handshake method to make secure websocket work (wss)
	wsServer := websocket.Server{
		Handshake: func(config *websocket.Config, req *http.Request) error {
			return nil
		},
		Handler: handler,
	}
	return wsServer.Handler
}

func main() {
	// Parse command line args
	flag.Parse()

	// Start gameBroadcaster
	go gameBroadcaster.Broadcast(0)
	go playersBroadcaster.Broadcast(0)

	// Start websocket server
	wsServerMux := http.NewServeMux()
	wsServerMux.Handle("/ws/game", newWebsocketHandler(rpsWebsocketHandler))
	wsServerMux.Handle("/ws/players", newWebsocketHandler(currentPlayersWebsocketHandler))
	go func() {
		http.ListenAndServe(":5001", wsServerMux)
	}()

	// Start main server
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/play", rpsHandler)
	serverMux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./client/public"))))
	log.Println("Serving at localhost:5000 and localhost:5001" )
	log.Fatal(http.ListenAndServe(":5000", serverMux))
}
