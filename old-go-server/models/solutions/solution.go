package models

type Solution struct {
	ProblemName string `json:"problem_name" binding:"required"`
	ProblemLink string `json:"problem_link"`
	ProblemID   string `json:"problem_id"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description" binding:"required"`
	Language    string `json:"language" binding:"required"`
	Difficulty  string `json:"difficulty"`
	Notes       string `json:"notes"`
}
