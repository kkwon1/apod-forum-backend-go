package models

type CommentTree struct {
	CommentID    string        `json:"commentId"`
	Children     []CommentTree `json:"children"`
	CreateDate   string        `json:"createDate"`
	ModifiedDate string        `json:"modifiedDate"`
	Comment      string        `json:"comment"`
	Author       string        `json:"author"`
	IsDeleted    bool          `json:"isDeleted"`
	IsLeaf       bool          `json:"isLeaf"`
}
