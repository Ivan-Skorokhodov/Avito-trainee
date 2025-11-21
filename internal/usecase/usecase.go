package usecase

import "PRmanager/internal/repository"

type UsecaseInterface interface {
}

type UseCase struct {
	repo repository.RepositoryInterface
}

func NewUseCase(repo repository.RepositoryInterface) *UseCase {
	return &UseCase{repo: repo}
}
