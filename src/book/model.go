package book

import "time"

type Book struct {
	ID          int       `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Author      string    `json:"author,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       string    `json:"price,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
