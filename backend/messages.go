package main

type IncomingMessage struct {
	Type     string `json:"type"`
	RoomCode string `json:"roomCode,omitempty"`
	Guess    string `json:"guess,omitempty"`
	Mode     string `json:"mode,omitempty"`
	Word     string `json:"word,omitempty"` // user submitted word
	Role     string `json:"role,omitempty"` // user chosen role
}

type OutgoingMessage struct {
	Type           string   `json:"type"`
	Role           string   `json:"role,omitempty"`
	Word           string   `json:"word,omitempty"`
	Words          []string `json:"words,omitempty"` // the 3 choices
	Outcome        string   `json:"outcome,omitempty"`
	ConfirmedRoles []string `json:"confirmedRoles,omitempty"` // roles already claimed
	RoundEndsAt    int64    `json:"roundEndsAt,omitempty"`    // unix milliseconds
}
