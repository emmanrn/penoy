package main

import (
	"math/rand"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

const codeChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

type Player struct {
	Conn      *websocket.Conn
	Role      string
	Confirmed bool
}

type Room struct {
	mu sync.Mutex // this is a mutex meaning, this would make only one goroutine run at a time and
	// others wait
	Code     string
	Players  []*Player
	Mode     string // this determines if the game is "random" (random words are used) | "custom" (user enters their own word)
	Words    []string
	WordIdx  int
	Current  string // current string regardless of mode
	TimerGen int    // this is used to check if the timer is still relevant, so any left over timers just does nothing
}

var (
	roomsMu sync.Mutex           // a mutex for all rooms so adding/removing rooms would only make one goroutine run at a time
	rooms   = map[string]*Room{} // list of all rooms in the server
)

func generateRoomCode() string {
	var sb strings.Builder
	for i := 0; i < 5; i++ {
		sb.WriteByte(codeChars[rand.Intn(len(codeChars))])
	}

	return sb.String()
}

func createRoom(mode string) *Room {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	var code string
	for {
		code = generateRoomCode()
		if _, exists := rooms[code]; !exists {
			break
		}
	}

	r := &Room{
		Code:  code,
		Mode:  mode,
		Words: shuffleWords(),
	}

	rooms[code] = r
	return r
}

// takes in 'mode' now to determine if its custom or random mode
func joinRoom(code string) *Room {
	// blocks other goroutines and only one goroutine may run
	roomsMu.Lock()

	// basically at the end of the function unlock the goroutine so other goroutines that are
	// waiting can take the next run
	defer roomsMu.Unlock() // lmao zig defer

	if r, ok := rooms[code]; ok {
		return r
	}

	// no room is found
	return nil
}

func (r *Room) broadcast(msg OutgoingMessage) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range r.Players {
		p.Conn.WriteJSON(msg)
	}
}

// TODO: DELETE
// // used to check which roles are currently taken
// func (r *Room) takenRoles() []string {
// 	var taken []string
// 	for _, p := range r.Players {
// 		if p.Role != "" {
// 			taken = append(taken, p.Role)
// 		}
// 	}
// 	return taken
// }

// used to check which roles are taken
func (r *Room) confirmedRoles() []string {
	var confirmed []string

	for _, p := range r.Players {
		if p.Confirmed {
			confirmed = append(confirmed, p.Role)
		}
	}

	return confirmed
}

// swap roles helper function
func (r *Room) swapRoles() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range r.Players {
		if p.Role == "clue-giver" {
			p.Role = "guesser"
		} else if p.Role == "guesser" {
			p.Role = "clue-giver"
		}
	}
}

func (r *Room) removePlayer(conn *websocket.Conn) {
	r.mu.Lock()
	remaining := r.Players[:0]
	for _, p := range r.Players {
		if p.Conn != conn {
			remaining = append(remaining, p)
		}
	}

	r.Players = remaining

	if len(r.Players) == 1 {
		r.Players[0].Role = ""
		r.Players[0].Confirmed = false
	}

	isEmpty := len(r.Players) == 0
	r.mu.Unlock()

	if isEmpty {
		roomsMu.Lock()
		delete(rooms, r.Code)
		roomsMu.Unlock()
	} else {
		r.broadcast(OutgoingMessage{Type: "opponent_left"})
	}
}
