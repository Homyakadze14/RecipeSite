package entities

type JWTToken struct {
	Token string `json:"token" binding:"required"`
}

type JWTData struct {
	UserID interface{} `json:"user_id"`
}
