package usecase

import (
	"PRmanager/internal/models"
	"PRmanager/internal/repository"
)

type UsecaseInterface interface {
	AddTeam(dto *models.TeamDTO) error
}

type UseCase struct {
	repo repository.RepositoryInterface
}

func NewUseCase(repo repository.RepositoryInterface) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) AddTeam(dto *models.TeamDTO) error {
	return nil
}
