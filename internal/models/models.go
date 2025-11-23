package models

import "time"

type Team struct {
	TeamId      int
	TeamName    string
	TeamMembers []*User
}

type User struct {
	UserId   int
	SystemId string
	UserName string
	TeamId   int
	TeamName string
	IsActive bool
	Reviews  []*PullRequest
}

type PullRequest struct {
	PullRequestId     int
	SystemId          string
	PullRequestName   string
	AuthorId          int
	AuthorSystemId    string
	Status            string
	AssigneeReviewers []*User
	CreatedAt         time.Time
	MergedAt          time.Time
}
