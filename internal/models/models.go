package models

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int64       `bun:"id,pk,autoincrement" json:"id" form:"id"`
	Email         string      `bun:"email,unique" json:"email"`
	Password      string      `bun:"password" json:"password"`
	Name          string      `bun:"name" json:"name"`
	Birthday      time.Time   `bun:"birthday" json:"birthday"`
	Sex           bool        `bun:"sex" json:"sex"`
	Latitude      float64     `bun:"latitude" json:"latitude"`
	Longitude     float64     `bun:"longitude" json:"longitude"`
	ImgPath       string      `bun:"img_path" json:"img_path"`
	Reaction      []*Reaction `bun:"rel:has-many,join:id=user_id"`
}

type Reaction struct {
	bun.BaseModel `bun:"table:reactions,alias:r"`
	UserID        int64 `bun:",pk" json:"user_id"`
	ToID          int64 `bun:",pk" json:"to_id"`
	To            *User `bun:"rel:belongs-to,join:to_id=id"`
	Like          bool  `bun:"like" json:"like"`
}
