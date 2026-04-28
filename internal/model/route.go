package model

import "time"

type Route struct {
	BotID     int64     `json:"bot_id"`
	TargetURL string    `json:"target_url"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
