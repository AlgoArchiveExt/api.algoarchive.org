package githubutils

import (
	"fmt"
	models "main/models/solutions"
	"regexp"
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

				difficultyRegex := regexp.MustCompile(`(?s)<h3>(.*?)</h3>`)
				matches := difficultyRegex.FindStringSubmatch(content)
				if len(matches) > 1 {
					solution.Difficulty = strings.ToLower(matches[1])
				}

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

// Returns a language string to its proper file extension (for example, "python" -> "py").
func MapLanguageStringToFileExtension(language string) (extension string, ok bool) {
	lowercaseString := strings.ToLower(language)

	switch lowercaseString {
	// Case for language strings that match their file extension
	case "cpp", "scala", "java", "c", "swift", "dart", "go", "php":
		return language, true
	case "python", "python3":
		return "py", true
	case "csharp":
		return "cs", true
	case "javascript":
		return "js", true
	case "typescript":
		return "ts", true
	case "kotlin":
		return "kt", true
	case "ruby":
		return "rb", true
	case "rust":
		return "rs", true
	case "racket":
		return "rkt", true
	case "erlang":
		return "erl", true
	case "elixir":
		return "ex", true
	}

	return "", false
}
