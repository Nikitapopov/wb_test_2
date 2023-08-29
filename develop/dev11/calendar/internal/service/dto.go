package service

// DTO для добавления
type InsertEventDTO struct {
	UserId      int    `json:"user_id"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

// DTO для обновления
type UpdateEventDTO struct {
	ID          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

// DTO для удаления
type RemoveEventDTO struct {
	ID int `json:"id"`
}
