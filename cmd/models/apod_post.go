package models

type ApodPost struct {
	NasaApod Apod        `json:"nasaApod"`
	Comments CommentTree `json:"comments"`
}
