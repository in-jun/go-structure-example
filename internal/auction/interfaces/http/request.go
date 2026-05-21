package http

import "time"

type CreateRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartPrice  int64     `json:"start_price"`
	EndTime     time.Time `json:"end_time"`
}
