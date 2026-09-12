package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

const roundDuration = 60 * time.Second

func handleWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer ws.Close()

	var room *Room

	defer func() {
		if room != nil {
			room.removePlayer(ws)
		}
	}()

	for {
		var msg IncomingMessage
		if err := ws.ReadJSON(&msg); err != nil {
			log.Println("read err: ", err)
			break
		}

		switch msg.Type {

		case "create_room":
			mode := msg.Mode
			if mode == "" {
				mode = "random"
			}

			room = createRoom(mode)

			room.mu.Lock()
			player := &Player{Conn: ws, Role: "", Confirmed: false}
			room.Players = append(room.Players, player)
			room.mu.Unlock()

			ws.WriteJSON(OutgoingMessage{Type: "room_joined", RoomCode: room.Code})

		case "join_room":
			roomsMu.Lock()
			existing, exists := rooms[msg.RoomCode]
			roomsMu.Unlock()

			if !exists {
				ws.WriteJSON(OutgoingMessage{Type: "room_not_found"})
				break
			}

			room = existing

			room.mu.Lock()
			player := &Player{Conn: ws, Role: "", Confirmed: false}
			room.Players = append(room.Players, player)
			full := len(room.Players) == 2
			room.mu.Unlock()

			ws.WriteJSON(OutgoingMessage{Type: "room_joined", RoomCode: room.Code})

			if full {
				room.broadcast(OutgoingMessage{Type: "choose_role_phase"})
			}

		case "submit_word":
			if room != nil {
				room.mu.Lock()
				mode := room.Mode
				room.mu.Unlock()
				if mode != "custom" {
					break
				}

				room.mu.Lock()
				room.Current = msg.Word
				room.mu.Unlock()

				endsAt := startRoundTimer(room)

				ws.WriteJSON(OutgoingMessage{Type: "word_locked", Word: msg.Word})
				for _, p := range room.Players {
					if p.Role == "guesser" {
						p.Conn.WriteJSON(OutgoingMessage{Type: "round_start", Role: p.Role, RoundEndsAt: endsAt})
					}
				}
			}
		case "next_word":
			if room != nil {
				nextWord(room)
			}
		case "leave_room":
			if room != nil {
				room.removePlayer(ws)
				room = nil
			}
		case "choose_role":
			if room != nil {
				room.mu.Lock()

				for _, p := range room.Players {
					if p.Conn == ws && !p.Confirmed {
						p.Role = msg.Role
					}
				}

				confirmed := room.confirmedRoles()
				room.mu.Unlock()

				room.broadcast(OutgoingMessage{Type: "role_state", ConfirmedRoles: confirmed})

			}

		case "cancel_role":
			if room != nil {
				room.mu.Lock()
				for _, p := range room.Players {
					if p.Conn == ws {
						p.Role = ""
						p.Confirmed = false
					}
				}

				confirmed := room.confirmedRoles()
				room.mu.Unlock()

				room.broadcast(OutgoingMessage{Type: "role_state", ConfirmedRoles: confirmed})
			}

		case "confirm_role":
			if room != nil {
				room.mu.Lock()
				var self *Player
				conflict := false

				for _, p := range room.Players {
					if p.Conn == ws {
						self = p
					} else if p.Confirmed && p.Role == msg.Role {
						conflict = true
					}
				}

				if self != nil && self.Role != "" && !conflict {
					self.Confirmed = true
				}

				confirm := room.confirmedRoles()
				bothConfirmed := len(room.Players) == 2 &&
					room.Players[0].Confirmed && room.Players[1].Confirmed &&
					room.Players[0].Role != room.Players[1].Role

				room.mu.Unlock()

				if conflict {
					ws.WriteJSON(OutgoingMessage{Type: "confirm_rejected", ConfirmedRoles: confirm})
				} else {
					room.broadcast(OutgoingMessage{Type: "role_state", ConfirmedRoles: confirm})
					if bothConfirmed {
						startRound(room)
					}
				}
			}

		case "choose_word":
			if room != nil {
				room.mu.Lock()
				room.Current = msg.Word
				room.WordIdx += 3
				room.mu.Unlock()

				endsAt := startRoundTimer(room)

				for _, p := range room.Players {
					if p.Role == "clue-giver" {
						p.Conn.WriteJSON(OutgoingMessage{Type: "round_start", Role: p.Role, Word: msg.Word, RoundEndsAt: endsAt})
					} else {
						p.Conn.WriteJSON(OutgoingMessage{Type: "round_start", Role: p.Role, RoundEndsAt: endsAt})
					}
				}
			}

		case "guess":
			if room != nil {
				currentRoom := room
				currentRoom.mu.Lock()
				current := strings.TrimSpace(strings.ToLower(room.Current))
				guess := strings.TrimSpace(strings.ToLower(msg.Guess))
				isCorrect := current != "" && current == guess

				var revealWord string
				var gen int

				if isCorrect {
					currentRoom.TimerGen++
					gen = currentRoom.TimerGen
					revealWord = currentRoom.Current
				}

				currentRoom.mu.Unlock()

				if isCorrect {
					currentRoom.broadcast(OutgoingMessage{Type: "correct_guess", Word: revealWord})

					time.AfterFunc(5*time.Second, func() {
						currentRoom.mu.Lock()
						stillCurrent := currentRoom.TimerGen == gen
						currentRoom.mu.Unlock()
						if stillCurrent {
							nextWord(room)
						}
					})
				}
			}

		}
	}

	log.Println("Client Successfully connected")
}

func startRound(room *Room) {
	room.mu.Lock()
	mode := room.Mode
	room.mu.Unlock()

	if mode == "custom" {
		for _, p := range room.Players {
			if p.Role == "clue-giver" {
				p.Conn.WriteJSON(OutgoingMessage{Type: "awaiting_word"})
			}
		}
		return
	}

	// random mode
	room.mu.Lock()
	if len(room.Words)-room.WordIdx < 3 {
		room.WordIdx = 0
		room.Words = shuffleWords()
	}
	choices := room.Words[room.WordIdx : room.WordIdx+3]
	room.mu.Unlock()

	for _, p := range room.Players {
		if p.Role == "clue-giver" {
			p.Conn.WriteJSON(OutgoingMessage{Type: "word_choices", Words: choices, Role: p.Role})
		} else if p.Role == "guesser" {
			p.Conn.WriteJSON(OutgoingMessage{Type: "waiting_word", Role: p.Role})
		}
	}
}

// returning an int64 here so we explicitlty return a 64-bit integer because unix milli is huge
// plus UnixMilli() already returns a int64 so we are just matching anyways
func startRoundTimer(room *Room) int64 {
	room.mu.Lock()
	room.TimerGen++
	gen := room.TimerGen
	room.mu.Unlock()

	endsAt := time.Now().Add(roundDuration)

	time.AfterFunc(roundDuration, func() {
		room.mu.Lock()
		stillCurrent := room.TimerGen == gen
		word := room.Current
		room.mu.Unlock()

		if stillCurrent {
			room.broadcast(OutgoingMessage{Type: "time_up", Word: word})
			room.TimerGen++
			time.AfterFunc(5*time.Second, func() {
				nextWord(room)
			})
		}
	})

	// we return milliseconds here for compatibility with the frontend
	// because JS native time unit is like in milliseconds??? so we just match that as well
	return endsAt.UnixMilli()
}

func nextWord(room *Room) {
	room.mu.Lock()
	room.TimerGen++
	room.Current = ""
	room.mu.Unlock()
	room.broadcast(OutgoingMessage{Type: "round_end"})
	room.swapRoles()
	startRound(room)
}
