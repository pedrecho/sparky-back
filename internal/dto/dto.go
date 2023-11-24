package dto

type Reaction struct {
	UserID int64  `json:"user_id"`
	ToID   int64  `json:"to_id"`
	Like   string `json:"like"`
}
