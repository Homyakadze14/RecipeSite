package models

import (
	"time"
)

type User struct {
	ID         int       `json:"id"`
	Email      string    `json:"email" validate:"required,email"`
	Login      string    `json:"login" validate:"required,min=3,max=20"`
	Password   string    `json:"password" validate:"required,min=8,max=50"`
	Icon_URL   string    `json:"icon_url"`
	About      string    `json:"about" validate:"max=1500"`
	Created_at time.Time `json:"created_at"`
}

type UserInfo struct {
	ID            int       `json:"id"`
	Login         string    `json:"login"`
	Icon_URL      string    `json:"icon_url"`
	About         string    `json:"about"`
	Created_at    time.Time `json:"created_at"`
	Recipies      []Recipe  `json:"recipies"`
	LikedRecipies []Recipe  `json:"liked_recipies"`
}

type Author struct {
	Login    string `json:"login"`
	Icon_URL string `json:"icon_url"`
}

type UserUpdate struct {
	Email    string `json:"email" validate:"email"`
	Login    string `json:"login" validate:"min=3,max=20"`
	Icon_URL string `json:"icon_url"`
	About    string `json:"about" validate:"max=1500"`
}

type UserPasswordUpdate struct {
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type UserLogin struct {
	Email    string `json:"email,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty" validate:"required"`
}
