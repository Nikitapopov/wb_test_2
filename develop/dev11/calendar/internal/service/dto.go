package service

type InsertEventDTO struct {
	UserId      int    `json:"user_id"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

type UpdateEventDTO struct {
	ID          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

type RemoveEventDTO struct {
	ID int `json:"id"`
}
