package usecase

import (
	"PRmanager/internal/models"
	"PRmanager/internal/repository"
	appErrors "PRmanager/pkg/app_errors"
	"PRmanager/pkg/logs"
	"context"
	"fmt"
	"math/rand/v2"
	"time"
)

type UsecaseInterface interface {
	AddTeam(ctx context.Context, dto *models.TeamDTO) error
	GetTeamByName(ctx context.Context, teamName string) (*models.TeamDTO, error)
	SetIsActive(ctx context.Context, dto *models.SetIsActiveDTO) (*models.UserDTO, error)
	GetReview(ctx context.Context, userSystemId string) (*models.ReviewDTO, error)
	CreatePullRequest(ctx context.Context, dto *models.InputCreatePullRequestDTO) (*models.OutputCreatePullRequestDTO, error)
	MergePullRequest(ctx context.Context, dto *models.InputMergePullRequestDTO) (*models.OutputMergePullRequestDTO, error)
	Reassign(ctx context.Context, dto *models.InputReassignDTO) (*models.OutputReassignDTO, error)
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

func (u *UseCase) CreatePullRequest(ctx context.Context, dto *models.InputCreatePullRequestDTO) (*models.OutputCreatePullRequestDTO, error) {
	exists, err := u.repo.PullRequestExists(ctx, dto.PullRequestId)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] CreatePullRequest", err.Error())
		return nil, appErrors.ErrServerError
	}

	if exists {
		logs.PrintLog(ctx, "[usecase] CreatePullRequest", appErrors.ErrPullRequestExists.Error())
		return nil, appErrors.ErrPullRequestExists
	}

	user, err := u.repo.GetUserBySystemId(ctx, dto.AuthorId)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] GetReview", err.Error())
		return nil, appErrors.ErrServerError
	}

	if user == nil {
		logs.PrintLog(ctx, "[usecase] GetReview", appErrors.ErrResourceNotFound.Error())
		return nil, appErrors.ErrResourceNotFound
	}

	members, err := u.repo.GetTeamMembers(ctx, user.TeamId)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] GetReview", err.Error())
		return nil, appErrors.ErrServerError
	}

	var candidates []*models.User

	for _, member := range members {
		if member.SystemId == dto.AuthorId {
			continue
		}
		if member.IsActive {
			candidates = append(candidates, member)
		}
	}

	var reviewers []*models.User
	if len(candidates) == 0 {
		// no one to review
	} else if len(candidates) == 1 {
		reviewers = append(reviewers, candidates[0])
	} else {
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})

		reviewers = append(reviewers, candidates[0])
		reviewers = append(reviewers, candidates[1])
	}
	logs.PrintLog(ctx, "[usecase] CreatePullRequest", fmt.Sprintf("Reviewers: %+v", reviewers))

	pr := &models.PullRequest{
		SystemId:        dto.PullRequestId,
		PullRequestName: dto.PullRequestName,
		AuthorId:        user.UserId,
		AuthorSystemId:  user.SystemId,
		Status:          "OPEN",
	}

	err = u.repo.CreatePullRequestAndReview(ctx, pr, reviewers)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] CreatePullRequest", err.Error())
		return nil, appErrors.ErrServerError
	}

	prDto := &models.OutputCreatePullRequestDTO{
		PullRequestID:     pr.SystemId,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorSystemId,
		Status:            pr.Status,
		AssignedReviewers: make([]string, 0, len(reviewers)),
	}

	for _, reviewer := range reviewers {
		prDto.AssignedReviewers = append(prDto.AssignedReviewers, reviewer.SystemId)
	}

	return prDto, nil
}

func (u *UseCase) MergePullRequest(ctx context.Context, dto *models.InputMergePullRequestDTO) (*models.OutputMergePullRequestDTO, error) {
	pr, err := u.repo.GetPullRequestById(ctx, dto.PullRequestId)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] MergePullRequest", err.Error())
		return nil, appErrors.ErrServerError
	}

	if pr == nil {
		logs.PrintLog(ctx, "[usecase] MergePullRequest", appErrors.ErrResourceNotFound.Error())
		return nil, appErrors.ErrResourceNotFound
	}

	if pr.Status == "MERGED" {
		logs.PrintLog(ctx, "[usecase] MergePullRequest", fmt.Sprintf("Pull request is already merged: name %+v id %+v", dto.PullRequestId, pr.PullRequestId))
		prDto := &models.OutputMergePullRequestDTO{
			PullRequestID:     pr.SystemId,
			PullRequestName:   pr.PullRequestName,
			AuthorID:          pr.AuthorSystemId,
			Status:            pr.Status,
			AssignedReviewers: make([]string, 0, len(pr.AssigneeReviewers)),
		}

		prDto.MergedAt = pr.MergedAt.Time.Format(time.RFC3339)

		for _, r := range pr.AssigneeReviewers {
			prDto.AssignedReviewers = append(prDto.AssignedReviewers, r.SystemId)
		}

		return prDto, nil
	}

	pr.Status = "MERGED"
	logs.PrintLog(ctx, "[usecase] MergePullRequest", fmt.Sprintf("Pull request is merged first time: name %+v id %+v", dto.PullRequestId, pr.PullRequestId))

	mergedTime, err := u.repo.SetMergedStatusPullRequest(ctx, pr.PullRequestId)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] MergePullRequest", err.Error())
		return nil, appErrors.ErrServerError
	}

	prDto := &models.OutputMergePullRequestDTO{
		PullRequestID:     pr.SystemId,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorSystemId,
		Status:            pr.Status,
		AssignedReviewers: make([]string, 0, len(pr.AssigneeReviewers)),
	}

	prDto.MergedAt = mergedTime.Time.Format(time.RFC3339)

	for _, r := range pr.AssigneeReviewers {
		prDto.AssignedReviewers = append(prDto.AssignedReviewers, r.SystemId)
	}

	return prDto, nil
}

func (u *UseCase) Reassign(ctx context.Context, dto *models.InputReassignDTO) (*models.OutputReassignDTO, error) {
	pr, err := u.repo.GetPullRequestById(ctx, dto.PullRequestId)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] MergePullRequest", err.Error())
		return nil, appErrors.ErrServerError
	}

	if pr == nil {
		logs.PrintLog(ctx, "[usecase] MergePullRequest", appErrors.ErrResourceNotFound.Error())
		return nil, appErrors.ErrResourceNotFound
	}

	user, err := u.repo.GetUserBySystemId(ctx, dto.UserId)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] GetReview", err.Error())
		return nil, appErrors.ErrServerError
	}

	if user == nil {
		logs.PrintLog(ctx, "[usecase] GetReview", appErrors.ErrResourceNotFound.Error())
		return nil, appErrors.ErrResourceNotFound
	}

	if pr.Status == "MERGED" {
		logs.PrintLog(ctx, "[usecase] Reassign", fmt.Sprintf("Pull request is already merged: name %+v id %+v", dto.PullRequestId, pr.PullRequestId))
		return nil, appErrors.ErrPullRequestMerged
	}

	IsUserReviewThisPR := false
	var otherReviewer *models.User
	for _, r := range pr.AssigneeReviewers {
		if r.SystemId == dto.UserId {
			IsUserReviewThisPR = true
			continue
		}
		otherReviewer = r
	}

	if !IsUserReviewThisPR {
		// return pr without replace
		prDto := &models.OutputReassignDTO{
			PullRequestID:     pr.SystemId,
			PullRequestName:   pr.PullRequestName,
			AuthorID:          pr.AuthorSystemId,
			Status:            pr.Status,
			AssignedReviewers: make([]string, 0, len(pr.AssigneeReviewers)),
			ReplacedBy:        "-",
		}

		for _, r := range pr.AssigneeReviewers {
			prDto.AssignedReviewers = append(prDto.AssignedReviewers, r.SystemId)
		}

		logs.PrintLog(ctx, "[usecase] Reassign", fmt.Sprintf("User is not assigned to pull request: name %+v id %+v", dto.PullRequestId, pr.PullRequestId))
		return prDto, nil
	}

	members, err := u.repo.GetTeamMembers(ctx, user.TeamId)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] GetReview", err.Error())
		return nil, appErrors.ErrServerError
	}

	var candidates []*models.User
	for _, member := range members {
		if member.SystemId == pr.AuthorSystemId || member.SystemId == dto.UserId || member.SystemId == otherReviewer.SystemId {
			continue
		}
		if member.IsActive {
			candidates = append(candidates, member)
		}
	}

	if len(candidates) == 0 {
		// just delete user from reviewers
	} else {
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})
	}

	err = u.repo.ReplaceReviewers(ctx, pr.PullRequestId, user.UserId, candidates[0].UserId)
	if err != nil {
		logs.PrintLog(ctx, "[usecase] Reassign", err.Error())
		return nil, appErrors.ErrServerError
	}

	return nil, nil
}
