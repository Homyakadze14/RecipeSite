package recipe

import "time"

type Recipe struct {
	ID          int       `json:"id"`
	User_ID     int       `json:"-"`
	Title       string    `json:"title" validate:"required,min=3,max=50"`
	About       string    `json:"about" validate:"required,max=2500"`
	Complexitiy int       `json:"complexitiy" validate:"required,min=1,max=3"`
	NeedTime    string    `json:"need_time" validate:"required"`
	Ingridients string    `json:"ingridients" validate:"required,max=1500"`
	Photos_URLS string    `json:"photos_urls"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}
