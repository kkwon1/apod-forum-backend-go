package models

type CommentNode struct {
	CommentID string         `json:"commentId"`
	Children  []*CommentNode `json:"children"`
	Comment   string         `json:"comment"`
	Author    string         `json:"author"`
}
