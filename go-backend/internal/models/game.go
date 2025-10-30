package models

import "time"

type Move struct {
	Player string    `bson:"player" json:"player"` // "R" or "Y"
	Col    int       `bson:"col" json:"col"`
	Row    int       `bson:"row" json:"row"`
	At     time.Time `bson:"at" json:"at"`
}

type GameDoc struct {
	GameID     string        `bson:"gameId" json:"gameId"`
	Player1    string        `bson:"player1" json:"player1"`
	Player2    string        `bson:"player2" json:"player2"`
	Winner     string        `bson:"winner" json:"winner"` // username or "Draw" or "Forfeit:<winner>"
	Duration   int           `bson:"duration" json:"duration"` // seconds
	FinalBoard [][]*string   `bson:"finalBoard" json:"finalBoard"`
	Moves      []Move        `bson:"moves" json:"moves"`
	CreatedAt  time.Time     `bson:"createdAt" json:"createdAt"`
}
