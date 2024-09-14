package entities

import (
	"io"
	"time"
)

type Recipe struct {
	ID           int       `json:"id"`
	UserID       int       `json:"creator_user_id"`
	Title        string    `json:"title" binding:"required,min=3,max=50"`
	About        string    `json:"about" binding:"required,max=2500"`
	Complexitiy  int       `json:"complexitiy" binding:"required,min=1,max=3"  enums:"1,2,3"`
	NeedTime     string    `json:"need_time" binding:"required"`
	Ingridients  string    `json:"ingridients" binding:"required,max=1500"`
	Instructions string    `json:"instructions" binding:"required,max=2000"`
	PhotosUrls   string    `json:"photos_urls"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateRecipe struct {
	Title        string          `json:"title" binding:"required,min=3,max=50"  form:"title"`
	About        string          `json:"about" binding:"required,max=2500"  form:"about"`
	Complexitiy  int             `json:"complexitiy" binding:"required,min=1,max=3"  enums:"1,2,3" form:"complexitiy"`
	NeedTime     string          `json:"need_time" binding:"required"  form:"need_time"`
	Ingridients  string          `json:"ingridients" binding:"required,max=1500"  form:"ingridients"`
	Instructions string          `json:"instructions" binding:"required,max=2000" form:"instructions"`
	Photos       []io.ReadSeeker `json:"-"`
}

func (r *CreateRecipe) HavePhotos() bool {
	return len(r.Photos) != 0
}

func (r *CreateRecipe) ToRecipe() *Recipe {
	return &Recipe{
		Title:        r.Title,
		About:        r.About,
		Complexitiy:  r.Complexitiy,
		NeedTime:     r.NeedTime,
		Ingridients:  r.Ingridients,
		Instructions: r.Instructions,
	}
}

type UpdateRecipe struct {
	Title        string          `json:"title" binding:"omitempty,min=3,max=50"  form:"title"`
	About        string          `json:"about" binding:"omitempty,max=2500"  form:"about"`
	Complexitiy  int             `json:"complexitiy" binding:"omitempty,min=1,max=3" enums:"1,2,3" form:"complexitiy"`
	NeedTime     string          `json:"need_time" binding:"omitempty"  form:"need_time"`
	Ingridients  string          `json:"ingridients" binding:"omitempty,max=1500"  form:"ingridients"`
	Instructions string          `json:"instructions" binding:"required,max=2000" form:"instructions"`
	Photos       []io.ReadSeeker `json:"-"`
}

func (r *UpdateRecipe) HavePhotos() bool {
	return len(r.Photos) != 0
}

func (r *UpdateRecipe) UpdateValues(recipe *Recipe) {
	if r.Complexitiy != 0 {
		recipe.Complexitiy = r.Complexitiy
	}
	if r.Title != "" {
		recipe.Title = r.Title
	}
	if r.About != "" {
		recipe.About = r.About
	}
	if r.NeedTime != "" {
		recipe.NeedTime = r.NeedTime
	}
	if r.Ingridients != "" {
		recipe.Ingridients = r.Ingridients
	}
	if r.Instructions != "" {
		recipe.Instructions = r.Instructions
	}
}

type RecipeInfo struct {
	Info *FullRecipe `json:"info"`
}

type FullRecipe struct {
	Recipe     *Recipe   `json:"recipe"`
	Author     *Author   `json:"author"`
	LikesCount int       `json:"likes_count"`
	IsLiked    bool      `json:"is_liked"`
	Comments   []Comment `json:"comments"`
}

type RecipeFilter struct {
	Limit      int    `json:"limit" example:"25"`
	Offset     int    `json:"offset" example:"0"`
	Query      string `json:"query" example:"tasty food"`
	OrderField string `json:"order_field" example:"title"  enums:"title,about,ingridients,emtpy"`
	OrderBy    int    `json:"order_by" binding:"min=-1,max=1"  enums:"-1,0,1"`
}
