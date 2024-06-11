package models

type Event struct {
	UserID   string `json:"user_id"`
	TaskID   string `json:"task_id"`
	AuthorID string `json:"author_id"`
}
