package models

type User struct {
	Owner string `json:"owner" binding:"required"`
	Repo  string `json:"repo_name" binding:"required"`
}
