package comment

import (
	"time"

	"github.com/Homyakadze14/RecipeSite/internal/user"
)

type Comment struct {
	ID        int          `json:"id"`
	UserID    int          `json:"-"`
	RecipeID  int          `json:"-"`
	Author    *user.Author `json:"author"`
	Text      string       `json:"text" validate:"required,min=1,max=250"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type CommentUpdate struct {
	ID   int    `json:"id" validate:"required"`
	Text string `json:"text" validate:"required,min=1,max=250"`
}

type CommentDelete struct {
	ID int `json:"id" validate:"required"`
}
