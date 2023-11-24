package models

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int64      `bun:"id,pk,autoincrement" json:"id"`
	Email         string     `bun:"email,unique" json:"email"`
	Password      string     `bun:"password" json:"password"`
	Name          string     `bun:"name" json:"name"`
	Birthday      time.Time  `bun:"birthday" json:"birthday"`
	Sex           bool       `bun:"sex" json:"sex"`
	Latitude      float64    `bun:"latitude" json:"latitude"`
	Longitude     float64    `bun:"longitude" json:"longitude"`
	ImgPath       string     `bun:"img_path" json:"img_path"`
	Reactions     []Reaction `bun:"rel:has-many,join:id=user_id"`
	Chats         []Chat     `bun:"m2m:user_chats,join:User=Chat"`
}

type Reaction struct {
	bun.BaseModel `bun:"table:reactions,alias:r"`
	UserID        int64 `bun:",pk" json:"user_id"`
	ToID          int64 `bun:",pk" json:"to_id"`
	To            *User `bun:"rel:belongs-to,join:to_id=id"`
	Like          bool  `bun:"like,default:false" json:"like"`
}

type Chat struct {
	bun.BaseModel `bun:"table:chats,alias:c"`
	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	Messages      []Message `bun:"rel:has-many,join:id=user_id" json:"messages"`
}

type UserChat struct {
	bun.BaseModel `bun:"table:user_chats,alias:uc"`
	UserID        int64 `bun:",pk"`
	User          *User `bun:"rel:belongs-to,join:user_id=id"`
	ChatID        int64 `bun:",pk"`
	Chat          *Chat `bun:"rel:belongs-to,join:chat_id=id"`
}

type Message struct {
	bun.BaseModel `bun:"table:messages,alias:m"`
	UserID        int64     `bun:",pk" json:"user_id" json:"user_id"`
	ChatID        int64     `bun:",pk" json:"chat_id" json:"chat_id"`
	Chat          *Chat     `bun:"rel:belongs-to,join:chat_id=id"`
	Time          time.Time `bun:"time" json:"time"`
	Text          string    `bun:"text" json:"text"`
}
