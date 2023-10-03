package models

import "time"

type Upvote struct {
	PostID   string    `json:"postId"`
	UserSub  string    `json:"userSub"`
	DateTime time.Time `json:"dateTime"`
}
