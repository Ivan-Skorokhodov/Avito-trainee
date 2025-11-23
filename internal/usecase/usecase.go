package usecase

import (
	"PRmanager/internal/models"
	"PRmanager/internal/repository"
	appErrors "PRmanager/pkg/app_errors"
	"PRmanager/pkg/logs"
	"context"
	"fmt"
)

type UsecaseInterface interface {
	AddTeam(ctx context.Context, dto *models.TeamDTO) error
	GetTeamByName(ctx context.Context, teamName string) (*models.TeamDTO, error)
	SetIsActive(ctx context.Context, dto *models.SetIsActiveDTO) (*models.UserDTO, error)
	GetReview(ctx context.Context, userSystemId string) (*models.ReviewDTO, error)
	CreatePullRequest(ctx context.Context, dto *models.CreatePullRequestDTO) (*models.PullRequestDTO, error)
}

type UseCase struct {
	repo repository.RepositoryInterface
}

func NewUseCase(repo repository.RepositoryInterface) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) AddTeam(ctx context.Context, dto *models.TeamDTO) error {
	exists, err := u.repo.TeamExists(ctx, dto.TeamName)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] AddTeam", err.Error())
		return appErrors.ErrServerError
	}
	if exists {
		logs.PrintLog(ctx, "[usecase] AddTeam", appErrors.ErrTeamExists.Error())
		return appErrors.ErrTeamExists
	}

	team := &models.Team{
		TeamName:    dto.TeamName,
		TeamMembers: make([]*models.User, 0, len(dto.Members)),
	}

	for _, m := range dto.Members {
		user := &models.User{
			SystemId: m.UserID,
			UserName: m.Username,
			IsActive: m.IsActive,
		}

		team.TeamMembers = append(team.TeamMembers, user)
	}

	if err := u.repo.CreateTeam(ctx, team); err != nil {
		logs.PrintLog(ctx, "[usecase] AddTeam", err.Error())
		return appErrors.ErrServerError
	}

	logs.PrintLog(ctx, "[usecase] AddTeam", fmt.Sprintf("Team added: %+v", team.TeamName))
	return nil
}

func (u *UseCase) GetTeamByName(ctx context.Context, teamName string) (*models.TeamDTO, error) {
	team, err := u.repo.GetTeamByName(ctx, teamName)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] GetTeamByName", err.Error())
		return nil, appErrors.ErrServerError
	}

	if team == nil {
		logs.PrintLog(ctx, "[usecase] GetTeamByName", appErrors.ErrResourceNotFound.Error())
		return nil, appErrors.ErrResourceNotFound
	}

	logs.PrintLog(ctx, "[usecase] GetTeamByName", fmt.Sprintf("Team found: %+v", team.TeamName))

	teamDto := &models.TeamDTO{
		TeamName: team.TeamName,
		Members:  make([]models.MemberDTO, 0, len(team.TeamMembers)),
	}

	for _, m := range team.TeamMembers {
		memberDTO := models.MemberDTO{
			UserID:   m.SystemId,
			Username: m.UserName,
			IsActive: m.IsActive,
		}

		teamDto.Members = append(teamDto.Members, memberDTO)
	}
	return teamDto, nil
}

func (u *UseCase) SetIsActive(ctx context.Context, dto *models.SetIsActiveDTO) (*models.UserDTO, error) {
	user, err := u.repo.SetIsActive(ctx, dto.UserID, dto.IsActive)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] SetIsActive", err.Error())
		return nil, appErrors.ErrServerError
	}

	if user == nil {
		logs.PrintLog(ctx, "[usecase] SetIsActive", appErrors.ErrResourceNotFound.Error())
		return nil, appErrors.ErrResourceNotFound
	}

	userDto := &models.UserDTO{
		UserId:   user.SystemId,
		UserName: user.UserName,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}

	logs.PrintLog(ctx, "[usecase] SetIsActive", fmt.Sprintf("Member updated: %+v set isActive to: %+v", dto.UserID, dto.IsActive))
	return userDto, nil
}

func (u *UseCase) GetReview(ctx context.Context, userSystemId string) (*models.ReviewDTO, error) {
	user, err := u.repo.GetUserBySystemId(ctx, userSystemId)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] GetReview", err.Error())
		return nil, appErrors.ErrServerError
	}

	if user == nil {
		logs.PrintLog(ctx, "[usecase] GetReview", appErrors.ErrResourceNotFound.Error())
		return nil, appErrors.ErrResourceNotFound
	}

	reviews, err := u.repo.GetListReviewsByUserId(ctx, user.UserId)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] GetReview", err.Error())
		return nil, appErrors.ErrServerError
	}

	reviewDto := &models.ReviewDTO{
		UserId:      user.SystemId,
		PullRequest: make([]models.PullRequestShortDTO, 0, len(reviews)),
	}

	for _, pr := range reviews {
		reviewDto.PullRequest = append(reviewDto.PullRequest, models.PullRequestShortDTO{
			PullRequestId:   pr.SystemId,
			PullRequestName: pr.PullRequestName,
			AuthorId:        pr.AuthorSystemId,
			Status:          pr.Status,
		})
	}

	logs.PrintLog(ctx, "[usecase] GetReview", fmt.Sprintf("Member found: %+v", user.SystemId))
	return reviewDto, nil
}

func (u *UseCase) CreatePullRequest(ctx context.Context, dto *models.CreatePullRequestDTO) (*models.PullRequestDTO, error) {
	exists, err := u.repo.PullRequestExists(ctx, dto.PullRequestId)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] CreatePullRequest", err.Error())
		return nil, appErrors.ErrServerError
	}

	if exists {
		logs.PrintLog(ctx, "[usecase] CreatePullRequest", appErrors.ErrPullRequestExists.Error())
		return nil, appErrors.ErrPullRequestExists
	}

	return nil, nil
}
