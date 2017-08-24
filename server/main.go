package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	websocket "github.com/gorilla/websocket"
)

var (
	port     = flag.String("port", ":5000", "Service port")
	upgrader = websocket.Upgrader{} // use default options

	leftTaken, rightTaken            = false, false
	rockWins, paperWins, scissorWins = 0, 0, 0
	rockTies, paperTies, scissorTies = 0, 0, 0
)

type Game struct {
	Left  string
	Right string
}

type RpsResponse struct {
	LeftTaken  bool
	RightTaken bool
	RpsWins    []int // Rock, Paper, Scissor wins
	RpsTies    []int // Rock, Paper, Scissor ties
}

func rpsHandler(w http.ResponseWriter, r *http.Request) {
}

func rpsWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()

	// Send initial response
	response := &RpsResponse{
		LeftTaken:  leftTaken,
		RightTaken: rightTaken,
		RpsWins:    []int{rockWins, paperWins, scissorWins},
		RpsTies:    []int{rockTies, paperTies, scissorTies},
	}
	message, _ := json.Marshal(response)
	err = c.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("write:", err)
		return
	}

	// Send responses on change
	for {
		break
	}
}

func main() {
	// Parse command line args
	flag.Parse()

	http.HandleFunc("/rps", rpsHandler)
	http.HandleFunc("/websocket/rps", rpsWebsocket)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("../client/public"))))
	log.Println("Serving at localhost" + *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}
