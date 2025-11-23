package models

type TeamDTO struct {
	TeamName string      `json:"team_name"`
	Members  []MemberDTO `json:"members"`
}

type MemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveDTO struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserDTO struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type ReviewDTO struct {
	UserId      string                `json:"user_id"`
	PullRequest []PullRequestShortDTO `json:"pull_requests"`
}

type PullRequestShortDTO struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
	Status          string `json:"status"`
}

type InputCreatePullRequestDTO struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
}

type OutputCreatePullRequestDTO struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

type InputMergePullRequestDTO struct {
	PullRequestId string `json:"pull_request_id"`
}

type OutputMergePullRequestDTO struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	MergedAt          string   `json:"merged_at"`
}
