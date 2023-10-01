package models

type Apod struct {
	Copyright      string `json:"copyright"`
	PostID         string `json:"postId"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	MediaType      string `json:"mediaType"`
	ServiceVersion string `json:"serviceVersion"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	Hdurl          string `json:"hdurl"`
	UpvoteCount    int    `json:"upvoteCount"`
	SaveCount      int    `json:"saveCount"`
	CommentCount   int    `json:"commentCount"`
}
