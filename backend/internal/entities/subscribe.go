package entities

type SubscribeInfo struct {
	CreatorID    int
	SubscriberID int
}

type SubscribeCreator struct {
	ID int `json:"creator_id" binding:"required"`
}

type NewRecipeRMQMessage struct {
	CreatorID int
	RecipeID  int
}
