package user

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Role      string    `json:"role"`
	Plan      string    `json:"plan"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
