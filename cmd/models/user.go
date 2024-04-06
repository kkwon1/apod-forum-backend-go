package models

type User struct {
	UserSub           string   `json:"userSub"`
	UserName          string   `json:"username"`
	Email             string   `json:"email"`
	EmailVerified     bool     `json:"emailVerified"`
	ProfilePictureUrl string   `json:"profilePictureUrl"`
	UpvotedPostIds    []string `json:"upvotedPostIds"`
}
