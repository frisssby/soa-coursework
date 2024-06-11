package models

type TaskStat struct {
	TaskID string `json:"task_id"`
	Views  uint64 `json:"views"`
	Likes  uint64 `json:"likes"`
}

type LikesStat struct {
	TaskID   string `json:"task_id"`
	AuthorID string `json:"author_id"`
	Likes    uint64 `json:"likes"`
}

type ViewsStat struct {
	TaskID   string `json:"task_id"`
	AuthorID string `json:"author_id"`
	Views    uint64 `json:"views"`
}

type UserStat struct {
	UserID string `json:"user_id"`
	Likes  uint64 `json:"likes"`
}
