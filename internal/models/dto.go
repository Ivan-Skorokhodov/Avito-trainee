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
