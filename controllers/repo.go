package controllers

import (
	"fmt"
	"log"
	"net/http"

	"main/infra/logger"
	githubutils "main/infra/utils/github"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v64/github"

	forms "main/forms/repo"
)

type RepositoryController struct{}

var repoForm = new(forms.RepositoryForm)

/*
Commit a solution to a repository given the owner's name, repository, problem name, and code.
https://github.com/AlgoArchiveExt/commit-testing
*/
func (ctrl *RepositoryController) CommitProblemSolution(c *gin.Context) {
	var commitForm = new(forms.CommitForm)

	if err := c.ShouldBindBodyWithJSON(&commitForm); err != nil {
		message := repoForm.Commit(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to parse body: %s", message),
		})

		return
	}

	gh := githubutils.CreateNewGithubClientWithUserToken(commitForm.AccessToken)

	// Get branch reference to the user's repository
	ref, _, err := gh.Git.GetRef(c, commitForm.User.Owner, commitForm.User.Repo, "heads/main")
	if err != nil {
		logger.Fatalf("Failed to get branch ref: %v", err)
		return
	}

	var basePath string = commitForm.Solution.ProblemName

	entries := []*github.TreeEntry{
		{
			Path:    github.String(basePath + "/" + commitForm.Solution.ProblemName + ".md"),
			Type:    github.String("blob"),
			Content: github.String(commitForm.Solution.Code),
			Mode:    github.String("100644"),
		},
		{
			Path:    github.String(basePath + "/" + "NOTES.md"),
			Type:    github.String("blob"),
			Content: github.String(commitForm.Solution.Notes),
			Mode:    github.String("100644"),
		},
		{
			Path:    github.String(basePath + "/" + commitForm.Solution.ProblemName + "." + commitForm.Solution.Language),
			Type:    github.String("blob"),
			Content: github.String(commitForm.Solution.Code),
			Mode:    github.String("100644"),
		},
	}

	// Create a new commit with the tree
	latestCommit, _, err := gh.Repositories.GetCommit(c, commitForm.User.Owner, commitForm.User.Repo, *ref.Object.SHA, nil)
	if err != nil {
		logger.Fatalf("Failed to get parent commit: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if latestCommit == nil || latestCommit.SHA == nil {
		logger.Fatalf("Latest commit tree or its SHA is nil")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Latest commit tree or its SHA is nil",
		})
		return
	}

	tree, _, err := gh.Git.CreateTree(c, commitForm.User.Owner, commitForm.User.Repo, *latestCommit.Commit.Tree.SHA, entries)
	if err != nil {
		logger.Fatalf("Failed to create tree: %v", err)
		return
	}

	if tree == nil || tree.SHA == nil {
		logger.Fatalf("Latest commit tree or its SHA is nil")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Latest commit tree or its SHA is nil",
		})
		return
	}

	commit := &github.Commit{
		Message: github.String(githubutils.CreateCommitMessage(commitForm.Solution.ProblemName)),
		Tree:    tree,
		Parents: []*github.Commit{{SHA: latestCommit.SHA}},
	}

	newCommit, _, err := gh.Git.CreateCommit(c, commitForm.User.Owner, commitForm.User.Repo, commit, nil)
	if err != nil {
		logger.Fatalf("Failed to create commit: %v", err)
		return
	}

	// Update the reference to point to the new commit
	ref.Object.SHA = newCommit.SHA
	_, _, err = gh.Git.UpdateRef(c, commitForm.User.Owner, commitForm.User.Repo, ref, false)
	if err != nil {
		log.Fatalf("Failed to update ref: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Successfully created commit at " + *newCommit.HTMLURL,
		"response": newCommit,
	})
}
