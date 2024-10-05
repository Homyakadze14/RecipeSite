package entities

import (
	"io"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email" binding:"required,email"`
	Login     string    `json:"login" binding:"required,min=3,max=20"`
	Password  string    `json:"password" binding:"required,min=8,max=50"`
	IconURL   string    `json:"icon_url"`
	About     string    `json:"about" binding:"max=1500"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthUser struct {
	Login     string `json:"login"`
	SessionID string `json:"session_id"`
}

type JSONUserInfo struct {
	User *UserInfo `json:"user"`
}

type UserInfo struct {
	ID            int                `json:"id"`
	Login         string             `json:"login"`
	IconURL       string             `json:"icon_url"`
	About         string             `json:"about"`
	CreatedAt     time.Time          `json:"created_at"`
	Recipies      []RecipeWithAuthor `json:"recipies"`
	LikedRecipies []RecipeWithAuthor `json:"liked_recipies"`
}

type Author struct {
	Login   string `json:"login"`
	IconURL string `json:"icon_url"`
}

type UserUpdate struct {
	Login string        `json:"login" binding:"omitempty,min=3,max=20" form:"login"`
	About string        `json:"about" binding:"omitempty,max=1500" form:"about"`
	Icon  io.ReadSeeker `json:"-"`
}

type UserIcon struct {
	IconURL string `json:"icon_url"`
}

func (u *UserUpdate) UpdateValues(user *User) {
	if u.Login != "" {
		user.Login = u.Login
	}
	if u.About != "" {
		user.About = u.About
	}
}

type UserPasswordUpdate struct {
	Password string `json:"password" binding:"required,min=8,max=50" example:"testpassword"`
}

type UserLogin struct {
	Email    string `json:"email,omitempty"  example:"test@test.com"`
	Login    string `json:"login,omitempty" example:"testuser" minlenght:"3" maxlenght:"8"`
	Password string `json:"password,omitempty" binding:"required,min=8,max=50" example:"testpassword"`
}
