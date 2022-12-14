package model

import "database/sql"

type Post struct {
	Id       int64
	DogId    int64
	Url      string
	Title    string
	Location string
}

type PostModel struct {
	Id       sql.NullInt64
	DogId    sql.NullInt64
	Url      string
	Title    string
	Location string
}

type PostRequest struct {
	Id       string `json:"id"`
	DogId    string `json:"dog"`
	Image    string `json:"image"`
	Url      string `json:"url"`
	Title    string `json:"title"`
	Location string `json:"location"`
}

type PostResponse struct {
	Id       string `json:"id"`
	DogId    string `json:"dog"`
	Image    string `json:"image"`
	Url      string `json:"url"`
	Title    string `json:"title"`
	Location string `json:"location"`
}
