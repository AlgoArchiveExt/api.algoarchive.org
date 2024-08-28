package controllers

import (
	"fmt"
	"net/http"

	"main/infra/logger"
	githubutils "main/infra/utils/github"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v64/github"

	forms "main/forms/solutions"
	formutils "main/infra/utils/forms"
)

type SolutionsController struct{}

/*
Commit a solution to a repository given the owner's name, repository, problem name, and code.
https://github.com/AlgoArchiveExt/commit-testing
*/
func (ctrl *SolutionsController) CommitProblemSolution(c *gin.Context) {
	var commitForm = &forms.CommitForm{}

	if err := c.ShouldBindBodyWithJSON(commitForm); err != nil {
		message := formutils.GenerateJSONBindingErrorMessage(commitForm, err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to parse body: %s", message),
		})

		return
	}

	ownerName := commitForm.User.Owner
	repoName := commitForm.User.Repo

	gh := githubutils.CreateNewGithubClientWithUserToken(commitForm.AccessToken)

	userRepoMainBranchReference, _, err := gh.Git.GetRef(c, ownerName, repoName, "heads/main")
	if err != nil {
		logger.Errorf("Failed to get branch ref: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to get reference for repo's main branch: %s", err.Error()),
		})
		return
	}

	basePath := commitForm.Solution.ProblemName

	commitRefPointsTo := userRepoMainBranchReference.Object

	latestCommit, _, err := gh.Repositories.GetCommit(c, ownerName, repoName, *commitRefPointsTo.SHA, nil)
	if err != nil {
		logger.Errorf("Failed to get parent commit: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to get parent commit for %s: %s", *userRepoMainBranchReference.Ref, err.Error()),
		})
		return
	}

	entries := []*github.TreeEntry{
		// Description
		{
			Path:    github.String(basePath + "/" + commitForm.Solution.ProblemName + ".md"),
			Type:    github.String("blob"),
			Content: github.String(commitForm.Solution.Description),
			Mode:    github.String("100644"),
		},
		// Notes
		{
			Path:    github.String(basePath + "/" + "NOTES.md"),
			Type:    github.String("blob"),
			Content: github.String(commitForm.Solution.Notes),
			Mode:    github.String("100644"),
		},
		// Solution Code
		{
			Path:    github.String(basePath + "/" + commitForm.Solution.ProblemName + "." + commitForm.Solution.Language),
			Type:    github.String("blob"),
			Content: github.String(commitForm.Solution.Code),
			Mode:    github.String("100644"),
		},
	}

	commitTree, _, err := gh.Git.CreateTree(c, ownerName, repoName, *latestCommit.Commit.Tree.SHA, entries)
	if err != nil {
		logger.Errorf("Failed to create tree: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to create tree from commit %s: %s", latestCommit.Commit, err.Error()),
		})
		return
	}

	commit := &github.Commit{
		Message: github.String(githubutils.CreateCommitMessage(commitForm.Solution.ProblemName)),
		Tree:    commitTree,
		Parents: []*github.Commit{{SHA: latestCommit.SHA}},
	}

	newCommit, _, err := gh.Git.CreateCommit(c, ownerName, repoName, commit, nil)
	if err != nil {
		logger.Errorf("Failed to create commit: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to create commit: %s", err.Error()),
		})
		return
	}

	commitRefPointsTo.SHA = newCommit.SHA
	_, _, err = gh.Git.UpdateRef(c, ownerName, repoName, userRepoMainBranchReference, false)
	if err != nil {
		logger.Errorf("Failed to update ref: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to update reference: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Successfully created commit at " + *newCommit.HTMLURL,
		"response": newCommit,
	})
}
