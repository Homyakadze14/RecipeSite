package entities

type SubscribeInfo struct {
	CreatorLogin string
	CreatorID    int
	SubscriberID int
}

type SubscribeCreator struct {
	ID int `json:"creator_id" binding:"required"`
}

type RecipeCreationMsg struct {
	CreatorID int
	RecipeID  int
}
