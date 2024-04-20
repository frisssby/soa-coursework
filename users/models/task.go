package models

type Task struct {
	UserId      string `json:"user_id,omitempty"`
	TaskId      string `json:"task_id,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}
