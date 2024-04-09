package models

type User struct {
	ID      	string `json:"id"`
	Email   	string `json:"email"`
	Name 		string `json:"name"`
	CreatedAt 	string `json:"created_at"`
	UpdatedAt 	string `json:"updated_at"`
	Role		string `json:"role"`
}