package models

import "time"

type Comment struct {
	PostID    string    `json:"postId"`
	CommentID string    `json:"commentId"`
	ParentID  string    `json:"parentId"`
	Comment   string    `json:"comment"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
