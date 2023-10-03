package requests

type UpvotePostRequest struct {
	PostId  string `json:"postId"`
	UserSub string `json:"userSub"`
}