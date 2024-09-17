package controllers

import (
	"fmt"
	"strings"

	"main/infra/logger"
	githubutils "main/infra/utils/github"
	"main/infra/utils/responses"
	models "main/models/solutions"

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

		responses.GiveErrorResponse(c, "Failed to parse body", message, nil)

		return
	}

	ownerName := commitForm.User.Owner
	repoName := commitForm.User.Repo

	gh := githubutils.CreateNewGithubClientWithUserToken(commitForm.AccessToken)

	userRepoMainBranchReference, _, err := gh.Git.GetRef(c, ownerName, repoName, "heads/main")
	if err != nil {
		logger.Errorf("Failed to get branch ref: %v", err)
		responses.GiveErrorResponse(c, "Failed to get reference for repo's main branch", err.Error(), nil)
		return
	}

	basePath := commitForm.Solution.ProblemName

	commitRefPointsTo := userRepoMainBranchReference.Object

	latestCommit, _, err := gh.Repositories.GetCommit(c, ownerName, repoName, *commitRefPointsTo.SHA, nil)
	if err != nil {
		logger.Errorf("Failed to get parent commit: %v", err)
		responses.GiveErrorResponse(c, fmt.Sprintf("Failed to get parent commit for %s", *userRepoMainBranchReference.Ref), err.Error(), nil)
		return
	}

	if commitForm.Solution.Difficulty != "" {
		// add difficulty to the top of the description <h3>difficulty</h3><hr>
		commitForm.Solution.Description = fmt.Sprintf("<h3>%s</h3><hr>%s", commitForm.Solution.Difficulty, commitForm.Solution.Description)
	}
	if commitForm.Solution.ProblemLink != "" && commitForm.Solution.ProblemID != "" {
		// add link to the problem at the bottom of the description <a href="problem_link">Problem Link</a>
		commitForm.Solution.Description = fmt.Sprintf("<h2><a href=\"%s\">%s. %s</a></h2>%s", commitForm.Solution.ProblemLink, commitForm.Solution.ProblemID, commitForm.Solution.ProblemName, commitForm.Solution.Description)
	}

	properFileExtension, ok := githubutils.MapLanguageStringToFileExtension(commitForm.Solution.Language)

	if !ok {
		responses.GiveErrorResponse(c, "Failed to parse language", "Could not parse language correctly, it might not be supported by us yet.", nil)
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
			Path:    github.String(basePath + "/" + githubutils.NotesFilename + ".md"),
			Type:    github.String("blob"),
			Content: github.String(commitForm.Solution.Notes),
			Mode:    github.String("100644"),
		},
		// Solution Code
		{
			Path:    github.String(basePath + "/" + commitForm.Solution.ProblemName + "." + properFileExtension),
			Type:    github.String("blob"),
			Content: github.String(commitForm.Solution.Code),
			Mode:    github.String("100644"),
		},
	}

	commitTree, _, err := gh.Git.CreateTree(c, ownerName, repoName, *latestCommit.Commit.Tree.SHA, entries)
	if err != nil {
		logger.Errorf("Failed to create tree: %v", err)
		responses.GiveErrorResponse(c, fmt.Sprintf("Failed to create tree from commit %s", *latestCommit.Commit.HTMLURL), err.Error(), nil)
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
		responses.GiveErrorResponse(c, "Failed to create commit", err.Error(), nil)
		return
	}

	commitRefPointsTo.SHA = newCommit.SHA
	_, _, err = gh.Git.UpdateRef(c, ownerName, repoName, userRepoMainBranchReference, false)
	if err != nil {
		logger.Errorf("Failed to update ref: %v", err)
		responses.GiveErrorResponse(c, "Failed to update reference", err.Error(), nil)
		return
	}

	responses.GiveOKResponse(c, "Successfully created commit at "+*newCommit.HTMLURL, &map[string]any{
		"response": newCommit,
	})
}

func (ctrl *SolutionsController) GetSolutions(c *gin.Context) {
	owner := c.Params.ByName("owner")
	repo := c.Params.ByName("repo")

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		responses.GiveUnauthorizedResponse(c, "Authorization header missing or invalid", nil)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	gh := githubutils.CreateNewGithubClientWithUserToken(token)

	mainBranchRef, _, err := gh.Git.GetRef(c, owner, repo, "heads/main")
	if err != nil {
		responses.GiveErrorResponse(c, fmt.Sprintf("Failed to get Git Reference for repo %s owned by user %s", repo, owner), err.Error(), nil)
		return
	}

	treeForLatestMainCommit, _, err := gh.Git.GetTree(c, owner, repo, *mainBranchRef.Object.SHA, false)
	if err != nil {
		responses.GiveErrorResponse(c, fmt.Sprintf("Failed to get Git Tree for repo %s owned by user %s", repo, owner), err.Error(), nil)
		return
	}

	/*
		As it stands, this function is very slow. It takes around 300 milliseconds for every solution to be parsed in the entire tree.
		I'll optimize it soon, but I want to get the server deployed so the front end team can start using this.

		There are a few optimizations I want to try:
		* Letting treeForLatestMainCommit be recursive so it gets every tree and blob from the latest commit and then sorting through those.
			* This may lead to performance gains because we would no longer have to get another tree from every problem, instead the github client
			  handles it for us.
		* Saving problem metadata in a database
			* We may be able to set up a connection to a SQL database and use raw SQL to look for all the problems and its
				metadata (including description, difficulty, topics, ect). This will let us save time getting every blob's raw byte data and information.
				Bonus points if the database is on the server itself.
	*/
	solutions := []models.Solution{}
	for _, entry := range treeForLatestMainCommit.Entries {
		if entry.GetType() == githubutils.Tree {
			var solution models.Solution = githubutils.ExtractSolutionFromTree(c, gh, owner, repo, entry)

			solutions = append(solutions, solution)
		}
	}

	responses.GiveOKResponse(c, fmt.Sprintf("Successfully obtained all solutions for repo %s!", owner+"/"+repo), &map[string]any{
		"solutions": solutions,
	})
}

func (ctrl *SolutionsController) GetSolutionsCount(c *gin.Context) {
	owner := c.Params.ByName("owner")
	repo := c.Params.ByName("repo")

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		responses.GiveUnauthorizedResponse(c, "Authorization header missing or invalid", nil)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	gh := githubutils.CreateNewGithubClientWithUserToken(token)

	mainBranchRef, _, err := gh.Git.GetRef(c, owner, repo, "heads/main")
	if err != nil {
		responses.GiveErrorResponse(c, fmt.Sprintf("Failed to get Git Reference for repo %s owned by user %s", repo, owner), err.Error(), nil)
		return
	}

	treeForLatestMainCommit, _, err := gh.Git.GetTree(c, owner, repo, *mainBranchRef.Object.SHA, false)
	if err != nil {
		responses.GiveErrorResponse(c, fmt.Sprintf("Failed to get Git Tree for repo %s owned by user %s", repo, owner), err.Error(), nil)
		return
	}

	solutionsCount := 0
	for _, entry := range treeForLatestMainCommit.Entries {
		if entry.GetType() == githubutils.Tree {
			solutionsCount++
		}
	}

	responses.GiveOKResponse(c, fmt.Sprintf("Successfully obtained all solutions for repo %s!", owner+"/"+repo), &map[string]any{
		"solutions_count": solutionsCount,
	})
}

func (ctrl *SolutionsController) GetSolutionsCountByDifficulty(c *gin.Context) {
	owner := c.Params.ByName("owner")
	repo := c.Params.ByName("repo")

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		responses.GiveUnauthorizedResponse(c, "Authorization header missing or invalid", nil)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	gh := githubutils.CreateNewGithubClientWithUserToken(token)

	mainBranchRef, _, err := gh.Git.GetRef(c, owner, repo, "heads/main")
	if err != nil {
		responses.GiveErrorResponse(c, fmt.Sprintf("Failed to get Git Reference for repo %s owned by user %s", repo, owner), err.Error(), nil)
		return
	}

	treeForLatestMainCommit, _, err := gh.Git.GetTree(c, owner, repo, *mainBranchRef.Object.SHA, false)
	if err != nil {
		responses.GiveErrorResponse(c, fmt.Sprintf("Failed to get Git Tree for repo %s owned by user %s", repo, owner), err.Error(), nil)
		return
	}

	difficulties := map[string]int{
		"easy":   0,
		"medium": 0,
		"hard":   0,
	}

	for _, entry := range treeForLatestMainCommit.Entries {
		if entry.GetType() == githubutils.Tree {
			solution := githubutils.ExtractSolutionFromTree(c, gh, owner, repo, entry)

			switch solution.Difficulty {
			case "easy":
				difficulties["easy"]++
			case "medium":
				difficulties["medium"]++
			case "hard":
				difficulties["hard"]++
			}
		}
	}

	responses.GiveOKResponse(c, fmt.Sprintf("Successfully obtained all solutions for repo %s!", owner+"/"+repo), &map[string]any{
		"difficulties": difficulties,
	})
}
