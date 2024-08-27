package githubutils

import (
	"fmt"

	"github.com/google/go-github/v64/github"
)

func CreateNewGithubClientWithUserToken(token string) *github.Client {
	client := github.NewClient(nil).WithAuthToken(token)

	return client
}

func CreateCommitMessage(problemName string) string {
	return fmt.Sprintf("COMMIT %s", problemName)
}
