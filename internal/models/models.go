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
	IsActive bool
}

type PullRequest struct {
	PullRequestId     string
	PullRequestName   string
	AuthorId          int
	Status            string
	AssigneeReviewers []*User
	CreatedAt         time.Time
	MergedAt          time.Time
}
