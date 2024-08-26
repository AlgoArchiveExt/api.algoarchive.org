package models

type User struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo_name"`
}
