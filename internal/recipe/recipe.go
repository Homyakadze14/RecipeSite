package recipe

import (
	"time"

	"github.com/Homyakadze14/RecipeSite/internal/comment"
	"github.com/Homyakadze14/RecipeSite/internal/user"
)

type Recipe struct {
	ID          int       `json:"id"`
	UserID      int       `json:"-"`
	Title       string    `json:"title" validate:"required,min=3,max=50"`
	About       string    `json:"about" validate:"required,max=2500"`
	Complexitiy int       `json:"complexitiy" validate:"required,min=1,max=3"`
	NeedTime    string    `json:"need_time" validate:"required"`
	Ingridients string    `json:"ingridients" validate:"required,max=1500"`
	PhotosUrls  string    `json:"photos_urls"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}

type FullRecipe struct {
	Recipe     *Recipe           `json:"recipe"`
	Author     *user.Author      `json:"author"`
	LikesCount int               `json:"likes_count"`
	IsLiked    bool              `json:"is_liked"`
	Comments   []comment.Comment `json:"comments"`
}

type RecipeFilter struct {
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	Query      string `json:"query"`
	OrderField string `json:"order_field"`
	OrderBy    int    `json:"order_by" validate:"min=-1,max=1"`
}
