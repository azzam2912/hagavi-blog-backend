package models

type BlogPost struct {
	ID      	string `json:"id"`
	Title   	string `json:"title"`
	Content 	string `json:"content"`
	CreatedAt 	string `json:"created_at"`
	UpdatedAt 	string `json:"updated_at"`
	Author		string `json:"author"`
}