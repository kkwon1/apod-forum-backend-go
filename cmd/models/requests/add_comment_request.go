package requests

type AddCommentRequest struct {
	PostID   string `json:"postId"`
	ParentID string `json:"parentId"`
	Comment  string `json:"comment"`
	Author   string `json:"author"`
}