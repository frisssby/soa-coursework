package models

type Task struct {
	UserID      string `json:"user_id,omitempty"`
	TaskID      string `json:"task_id,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}
