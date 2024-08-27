package forms

import (
	models "main/models/repo"
)

type CommitForm struct {
	AccessToken string          `json:"user_access_token" binding:"required"`
	User        models.User     `json:"user" binding:"required"`
	Solution    models.Solution `json:"solution" binding:"required"`
}
