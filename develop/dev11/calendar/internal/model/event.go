package model

const (
	DateLayout = "2006-01-02"
)

type Event struct {
	ID          int    `json:"id,omitempty"`
	UserId      int    `json:"user_id"`
	Date        string `json:"date"`
	RemoveDate  string `json:"remove_date,omitempty"`
	Description string `json:"description"`
}
