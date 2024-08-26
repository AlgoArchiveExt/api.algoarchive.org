package forms

import (
	"encoding/json"
	"main/infra/logger"
	models "main/models/repo"

	"github.com/go-playground/validator/v10"
)

type RepositoryForm struct{}

type CommitForm struct {
	AccessToken string          `json:"user_access_token" binding:"required"`
	User        models.User     `json:"user" binding:"required"`
	Solution    models.Solution `json:"solution" binding:"required"`
}

func (form RepositoryForm) AccessToken(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		return "User Access Token is required"
	default:
		return "User Access Token is invalid"
	}
}

func (form RepositoryForm) User(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		return "User is required"
	default:
		return "User is invalid"
	}
}

func (form RepositoryForm) Solution(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		return "Solution is required"
	default:
		return "Solution is invalid"
	}
}

func (form RepositoryForm) Commit(err error) string {
	logger.Errorf("ran into error: %s of type %T", err, err)

	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Some fields are invalid or missing"
		}

		for _, e := range err.(validator.ValidationErrors) {
			// logger.Infof("formatting %s", e.Field())
			switch e.Field() {
			case "AccessToken":
				return form.AccessToken(e.Tag())
			case "User":
				return form.User(e.Tag())
			case "Solution":
				return form.Solution(e.Tag())
			}
		}
	case *json.SyntaxError:
		return "Syntax Error"

	default:
		return "Invalid request"
	}

	return "Some fields are invalid"
}
