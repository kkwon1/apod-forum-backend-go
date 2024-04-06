package models

import "time"

type Upvote struct {
	PostId   string    `json:"postId" bson:"postId"`
	UserSub  string    `json:"userSub" bson:"userSub"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}