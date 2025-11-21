package usecase

import (
	"PRmanager/internal/models"
	"PRmanager/internal/repository"
	appErrors "PRmanager/pkg/app_errors"
	"PRmanager/pkg/logs"
	"context"
)

type UsecaseInterface interface {
	AddTeam(ctx context.Context, dto *models.TeamDTO) error
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
		return appErrors.ServerError
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
		return appErrors.ServerError
	}

	return nil
}
