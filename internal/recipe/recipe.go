package recipe

import "time"

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
	Recipe        *Recipe `json:"recipe"`
	Author        string  `json:"author"`
	AuthorIconUrl string  `json:"author_icon_url"`
}

type RecipeFilter struct {
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	Query      string `json:"query"`
	OrderField string `json:"order_field"`
	OrderBy    int    `json:"order_by" validate:"min=-1,max=1"`
}
