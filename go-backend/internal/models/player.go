package models

type Player struct {
	Username string `bson:"username" json:"username"`
	Wins     int    `bson:"wins" json:"wins"`
	Losses   int    `bson:"losses" json:"losses"`
	Draws    int    `bson:"draws" json:"draws"`
}
