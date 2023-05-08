package review

import "time"

type Review struct {
	ID        int       `json:"id,omitempty"`
	UserID    int       `json:"user_id,omitempty"`
	BookID    int       `json:"book_id,omitempty"`
	Rating    int       `json:"rating,omitempty"`
	Comment   string    `json:"comment,omitempty"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
