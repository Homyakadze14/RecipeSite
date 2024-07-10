package user

import (
	"time"
)

type User struct {
	ID         int       `json:"id,omitempty"`
	Email      string    `json:"email,omitempty" validate:"required,email"`
	Login      string    `json:"login,omitempty" validate:"required,min=3,max=20"`
	Password   string    `json:"password,omitempty" validate:"required,min=8,max=50"`
	Icon_URL   string    `json:"icon_url,omitempty"`
	About      string    `json:"about,omitempty" validate:"required,max=1500"`
	Created_at time.Time `json:"created_at,omitempty"`
}

type UserLogin struct {
	Email    string `json:"email,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty" validate:"required"`
}
