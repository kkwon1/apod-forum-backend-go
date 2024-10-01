package models

import "time"

type Comment struct {
	PostID    string    `json:"postId" bson:"postId"`
	CommentID string    `json:"commentId" bson:"commentId"`
	ParentID  string    `json:"parentId" bson:"parentId"`
	Comment   string    `json:"comment" bson:"comment"`
	Author    string    `json:"author" bson:"author"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
