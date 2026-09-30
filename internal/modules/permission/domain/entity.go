package domain

import "time"

type Permission struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Module    string    `json:"module"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
