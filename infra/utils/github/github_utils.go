package githubutils

import (
	"fmt"
	models "main/models/solutions"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v64/github"
)

func CreateNewGithubClientWithUserToken(token string) *github.Client {
	client := github.NewClient(nil).WithAuthToken(token)

	return client
}

func CreateCommitMessage(problemName string) string {
	return fmt.Sprintf("COMMIT %s", problemName)
}

func ExtractSolutionFromTree(c *gin.Context, gh *github.Client, owner string, repo string, entry *github.TreeEntry) models.Solution {
	treeForThisProblem, _, _ := gh.Git.GetTree(c, owner, repo, *entry.SHA, true)

	problemName := entry.Path
	solution := models.Solution{
		ProblemName: *problemName,
	}

	for _, blob := range treeForThisProblem.Entries {
		rawBlobContent, _, _ := gh.Git.GetBlobRaw(c, owner, repo, *blob.SHA)
		content := string(rawBlobContent)

		pathSplit := strings.Split(*blob.Path, ".")
		filename, fileExtension := pathSplit[0], pathSplit[1]

		switch fileExtension {
		case "md":
			// Description file
			if filename == *problemName || filename == "README" {
				solution.Description = content
				// Notes file
			} else if filename == NotesFilename {
				solution.Notes = content
			}
		default:
			solution.Code = content
			solution.Language = fileExtension
		}
	}

	return solution
}

const (
	Blob          string = "blob"
	Tree          string = "tree"
	NotesFilename string = "NOTES"
)
