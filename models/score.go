package models

import "time"

type Score struct {
	ID         int       `json:"id"`
	PlayerName string    `json:"player_name"`
	Score      int       `json:"score"`
	CreatedAt  time.Time `json:"created_at"`
}
