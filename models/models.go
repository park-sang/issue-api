package models

import "time"

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Issue struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	User        *User     `json:"user,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

var Users = []User{
	{ID: 1, Name: "김개발"},
	{ID: 2, Name: "이디자인"},
	{ID: 3, Name: "박기획"},
}
