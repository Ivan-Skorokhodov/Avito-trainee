package usecase_test

import (
	"PRmanager/internal/models"
	"PRmanager/internal/repository/mocks"
	"PRmanager/internal/usecase"
	appErrors "PRmanager/pkg/app_errors"
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUseCase_AddTeam(t *testing.T) {
	tests := []struct {
		name        string
		dto         *models.TeamDTO
		mockSetup   func(m *mocks.MockRepositoryInterface)
		expectedErr error
	}{
		{
			name: "success create team",
			dto: &models.TeamDTO{
				TeamName: "backend",
				Members: []models.MemberDTO{
					{UserID: "u1", Username: "Nick", IsActive: true},
					{UserID: "u2", Username: "Sara", IsActive: false},
				},
			},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().TeamExists(gomock.Any(), "backend").Return(false, nil)
				m.EXPECT().
					CreateTeam(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, team *models.Team) error {
						assert.Equal(t, "backend", team.TeamName)
						assert.Len(t, team.TeamMembers, 2)
						assert.Equal(t, "u1", team.TeamMembers[0].SystemId)
						assert.Equal(t, "Nick", team.TeamMembers[0].UserName)
						return nil
					})
			},
			expectedErr: nil,
		},
		{
			name: "team exists",
			dto:  &models.TeamDTO{TeamName: "backend"},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().TeamExists(gomock.Any(), "backend").Return(true, nil)
			},
			expectedErr: appErrors.ErrTeamExists,
		},
		{
			name: "error checking team exists",
			dto:  &models.TeamDTO{TeamName: "backend"},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().TeamExists(gomock.Any(), "backend").Return(false, errors.New("db error"))
			},
			expectedErr: appErrors.ErrServerError,
		},
		{
			name: "error in CreateTeam",
			dto: &models.TeamDTO{
				TeamName: "backend",
				Members:  []models.MemberDTO{},
			},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().TeamExists(gomock.Any(), "backend").Return(false, nil)
				m.EXPECT().CreateTeam(gomock.Any(), gomock.Any()).Return(errors.New("insert error"))
			},
			expectedErr: appErrors.ErrServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepositoryInterface(ctrl)
			uc := usecase.NewUseCase(mockRepo)

			tt.mockSetup(mockRepo)

			err := uc.AddTeam(context.Background(), tt.dto)

			if tt.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			}
		})
	}
}

func TestUseCase_GetTeamByName(t *testing.T) {
	tests := []struct {
		name        string
		teamName    string
		mockSetup   func(m *mocks.MockRepositoryInterface)
		expected    *models.TeamDTO
		expectedErr error
	}{
		{
			name:     "success — team found",
			teamName: "backend",
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().
					GetTeamByName(gomock.Any(), "backend").
					Return(&models.Team{
						TeamName: "backend",
						TeamMembers: []*models.User{
							{SystemId: "u1", UserName: "Nick", IsActive: true},
							{SystemId: "u2", UserName: "Sara", IsActive: false},
						},
					}, nil)
			},
			expected: &models.TeamDTO{
				TeamName: "backend",
				Members: []models.MemberDTO{
					{UserID: "u1", Username: "Nick", IsActive: true},
					{UserID: "u2", Username: "Sara", IsActive: false},
				},
			},
			expectedErr: nil,
		},
		{
			name:     "error from repository",
			teamName: "backend",
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().
					GetTeamByName(gomock.Any(), "backend").
					Return(nil, errors.New("db error"))
			},
			expected:    nil,
			expectedErr: appErrors.ErrServerError,
		},
		{
			name:     "team not found (nil result)",
			teamName: "backend",
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().
					GetTeamByName(gomock.Any(), "backend").
					Return(nil, nil)
			},
			expected:    nil,
			expectedErr: appErrors.ErrResourceNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepositoryInterface(ctrl)
			uc := usecase.NewUseCase(mockRepo)

			tt.mockSetup(mockRepo)

			result, err := uc.GetTeamByName(context.Background(), tt.teamName)

			if tt.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.TeamName, result.TeamName)
				assert.Equal(t, tt.expected.Members, result.Members)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Nil(t, result)
			}
		})
	}
}

func TestUseCase_SetIsActive(t *testing.T) {
	tests := []struct {
		name        string
		dto         *models.SetIsActiveDTO
		mockSetup   func(m *mocks.MockRepositoryInterface)
		expected    *models.UserDTO
		expectedErr error
	}{
		{
			name: "success — user updated",
			dto: &models.SetIsActiveDTO{
				UserID:   "u1",
				IsActive: true,
			},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().
					SetIsActive(gomock.Any(), "u1", true).
					Return(&models.User{
						SystemId: "u1",
						UserName: "Nick",
						TeamName: "backend",
						IsActive: true,
					}, nil)
			},
			expected: &models.UserDTO{
				UserId:   "u1",
				UserName: "Nick",
				TeamName: "backend",
				IsActive: true,
			},
			expectedErr: nil,
		},
		{
			name: "repo returns error",
			dto: &models.SetIsActiveDTO{
				UserID:   "u1",
				IsActive: false,
			},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().
					SetIsActive(gomock.Any(), "u1", false).
					Return(nil, errors.New("db error"))
			},
			expected:    nil,
			expectedErr: appErrors.ErrServerError,
		},
		{
			name: "user not found",
			dto: &models.SetIsActiveDTO{
				UserID:   "u1",
				IsActive: true,
			},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().
					SetIsActive(gomock.Any(), "u1", true).
					Return(nil, nil)
			},
			expected:    nil,
			expectedErr: appErrors.ErrResourceNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepositoryInterface(ctrl)
			uc := usecase.NewUseCase(mockRepo)

			tt.mockSetup(mockRepo)

			result, err := uc.SetIsActive(context.Background(), tt.dto)

			if tt.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Nil(t, result)
			}
		})
	}
}

func TestUseCase_GetReview(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		mockSetup   func(m *mocks.MockRepositoryInterface)
		expected    *models.ReviewDTO
		expectedErr error
	}{
		{
			name:   "success — user found with PRs",
			userID: "u1",
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				// user found
				m.EXPECT().
					GetUserBySystemId(gomock.Any(), "u1").
					Return(&models.User{
						UserId:   10,
						SystemId: "u1",
					}, nil)

				// PR reviews
				m.EXPECT().
					GetListReviewsByUserId(gomock.Any(), 10).
					Return([]*models.PullRequest{
						{
							SystemId:        "PR1",
							PullRequestName: "Fix Bug",
							AuthorSystemId:  "u2",
							Status:          "OPEN",
						},
						{
							SystemId:        "PR2",
							PullRequestName: "Refactor",
							AuthorSystemId:  "u3",
							Status:          "MERGED",
						},
					}, nil)
			},
			expected: &models.ReviewDTO{
				UserId: "u1",
				PullRequest: []models.PullRequestShortDTO{
					{
						PullRequestId:   "PR1",
						PullRequestName: "Fix Bug",
						AuthorId:        "u2",
						Status:          "OPEN",
					},
					{
						PullRequestId:   "PR2",
						PullRequestName: "Refactor",
						AuthorId:        "u3",
						Status:          "MERGED",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:   "repo error on GetUserBySystemId",
			userID: "u1",
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().
					GetUserBySystemId(gomock.Any(), "u1").
					Return(nil, errors.New("db fail"))
			},
			expected:    nil,
			expectedErr: appErrors.ErrServerError,
		},
		{
			name:   "user not found",
			userID: "u1",
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().
					GetUserBySystemId(gomock.Any(), "u1").
					Return(nil, nil)
			},
			expected:    nil,
			expectedErr: appErrors.ErrResourceNotFound,
		},
		{
			name:   "repo error on GetListReviewsByUserId",
			userID: "u1",
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				// first call: user exists
				m.EXPECT().
					GetUserBySystemId(gomock.Any(), "u1").
					Return(&models.User{
						UserId:   10,
						SystemId: "u1",
					}, nil)

				// second call fails
				m.EXPECT().
					GetListReviewsByUserId(gomock.Any(), 10).
					Return(nil, errors.New("db fail"))
			},
			expected:    nil,
			expectedErr: appErrors.ErrServerError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepositoryInterface(ctrl)
			uc := usecase.NewUseCase(mockRepo)

			tt.mockSetup(mockRepo)

			result, err := uc.GetReview(context.Background(), tt.userID)

			if tt.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.UserId, result.UserId)
				assert.Equal(t, tt.expected.PullRequest, result.PullRequest)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Nil(t, result)
			}
		})
	}
}

func TestUseCase_CreatePullRequest(t *testing.T) {
	tests := []struct {
		name      string
		dto       *models.InputCreatePullRequestDTO
		mockSetup func(m *mocks.MockRepositoryInterface)
		check     func(t *testing.T, out *models.OutputCreatePullRequestDTO, err error)
	}{
		{
			name: "error PullRequestExists",
			dto:  &models.InputCreatePullRequestDTO{PullRequestId: "PR1", AuthorId: "u1"},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().PullRequestExists(gomock.Any(), "PR1").Return(false, errors.New("db error"))
			},
			check: func(t *testing.T, out *models.OutputCreatePullRequestDTO, err error) {
				assert.Nil(t, out)
				assert.Equal(t, appErrors.ErrServerError, err)
			},
		},
		{
			name: "PR already exists",
			dto:  &models.InputCreatePullRequestDTO{PullRequestId: "PR1", AuthorId: "u1"},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().PullRequestExists(gomock.Any(), "PR1").Return(true, nil)
			},
			check: func(t *testing.T, out *models.OutputCreatePullRequestDTO, err error) {
				assert.Nil(t, out)
				assert.Equal(t, appErrors.ErrPullRequestExists, err)
			},
		},
		{
			name: "error GetUserBySystemId",
			dto:  &models.InputCreatePullRequestDTO{PullRequestId: "PR1", AuthorId: "u1"},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().PullRequestExists(gomock.Any(), "PR1").Return(false, nil)
				m.EXPECT().GetUserBySystemId(gomock.Any(), "u1").Return(nil, errors.New("db"))
			},
			check: func(t *testing.T, out *models.OutputCreatePullRequestDTO, err error) {
				assert.Nil(t, out)
				assert.Equal(t, appErrors.ErrServerError, err)
			},
		},
		{
			name: "user not found",
			dto:  &models.InputCreatePullRequestDTO{PullRequestId: "PR1", AuthorId: "u1"},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().PullRequestExists(gomock.Any(), "PR1").Return(false, nil)
				m.EXPECT().GetUserBySystemId(gomock.Any(), "u1").Return(nil, nil)
			},
			check: func(t *testing.T, out *models.OutputCreatePullRequestDTO, err error) {
				assert.Nil(t, out)
				assert.Equal(t, appErrors.ErrResourceNotFound, err)
			},
		},
		{
			name: "error GetTeamMembers",
			dto:  &models.InputCreatePullRequestDTO{PullRequestId: "PR1", AuthorId: "u1"},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().PullRequestExists(gomock.Any(), "PR1").Return(false, nil)
				m.EXPECT().GetUserBySystemId(gomock.Any(), "u1").
					Return(&models.User{UserId: 10, TeamId: 99, SystemId: "u1"}, nil)
				m.EXPECT().GetTeamMembers(gomock.Any(), 99).
					Return(nil, errors.New("db"))
			},
			check: func(t *testing.T, out *models.OutputCreatePullRequestDTO, err error) {
				assert.Nil(t, out)
				assert.Equal(t, appErrors.ErrServerError, err)
			},
		},
		{
			name: "0 reviewers",
			dto:  &models.InputCreatePullRequestDTO{PullRequestId: "PR1", PullRequestName: "Fix", AuthorId: "u1"},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().PullRequestExists(gomock.Any(), "PR1").Return(false, nil)
				m.EXPECT().GetUserBySystemId(gomock.Any(), "u1").
					Return(&models.User{UserId: 10, TeamId: 5, SystemId: "u1"}, nil)
				m.EXPECT().GetTeamMembers(gomock.Any(), 5).
					Return([]*models.User{{SystemId: "u1", IsActive: true}}, nil)
				m.EXPECT().CreatePullRequestAndReview(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ *models.PullRequest, reviewers []*models.User) error {
						assert.Len(t, reviewers, 0)
						return nil
					})
			},
			check: func(t *testing.T, out *models.OutputCreatePullRequestDTO, err error) {
				assert.NoError(t, err)
				assert.Len(t, out.AssignedReviewers, 0)
			},
		},
		{
			name: "1 reviewer",
			dto:  &models.InputCreatePullRequestDTO{PullRequestId: "PR1", AuthorId: "u1"},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().PullRequestExists(gomock.Any(), "PR1").Return(false, nil)
				m.EXPECT().GetUserBySystemId(gomock.Any(), "u1").
					Return(&models.User{UserId: 10, TeamId: 5, SystemId: "u1"}, nil)
				m.EXPECT().GetTeamMembers(gomock.Any(), 5).
					Return([]*models.User{
						{SystemId: "u1"},
						{SystemId: "u2", IsActive: true},
					}, nil)
				m.EXPECT().CreatePullRequestAndReview(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ *models.PullRequest, reviewers []*models.User) error {
						assert.Len(t, reviewers, 1)
						assert.Equal(t, "u2", reviewers[0].SystemId)
						return nil
					})
			},
			check: func(t *testing.T, out *models.OutputCreatePullRequestDTO, err error) {
				assert.NoError(t, err)
				assert.Len(t, out.AssignedReviewers, 1)
			},
		},
		{
			name: "2 reviewers",
			dto:  &models.InputCreatePullRequestDTO{PullRequestId: "PR1", AuthorId: "u1"},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().PullRequestExists(gomock.Any(), "PR1").Return(false, nil)

				m.EXPECT().GetUserBySystemId(gomock.Any(), "u1").
					Return(&models.User{UserId: 10, TeamId: 9, SystemId: "u1"}, nil)

				m.EXPECT().GetTeamMembers(gomock.Any(), 9).
					Return([]*models.User{
						{SystemId: "u1"},
						{SystemId: "u2", IsActive: true},
						{SystemId: "u3", IsActive: true},
					}, nil)

				m.EXPECT().
					CreatePullRequestAndReview(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ *models.PullRequest, reviewers []*models.User) error {
						assert.Len(t, reviewers, 2)
						return nil
					})
			},
			check: func(t *testing.T, out *models.OutputCreatePullRequestDTO, err error) {
				assert.NoError(t, err)
				assert.Len(t, out.AssignedReviewers, 2)
			},
		},
		{
			name: "error CreatePullRequestAndReview",
			dto:  &models.InputCreatePullRequestDTO{PullRequestId: "PR1", AuthorId: "u1"},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().PullRequestExists(gomock.Any(), "PR1").Return(false, nil)

				m.EXPECT().GetUserBySystemId(gomock.Any(), "u1").
					Return(&models.User{UserId: 10, TeamId: 9, SystemId: "u1"}, nil)

				m.EXPECT().GetTeamMembers(gomock.Any(), 9).
					Return([]*models.User{}, nil)

				m.EXPECT().CreatePullRequestAndReview(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("db"))
			},
			check: func(t *testing.T, out *models.OutputCreatePullRequestDTO, err error) {
				assert.Nil(t, out)
				assert.Equal(t, appErrors.ErrServerError, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepositoryInterface(ctrl)
			uc := usecase.NewUseCase(mockRepo)

			tt.mockSetup(mockRepo)

			out, err := uc.CreatePullRequest(context.Background(), tt.dto)
			tt.check(t, out, err)
		})
	}
}

func TestUseCase_MergePullRequest(t *testing.T) {
	tests := []struct {
		name      string
		dto       *models.InputMergePullRequestDTO
		mockSetup func(m *mocks.MockRepositoryInterface)
		check     func(t *testing.T, out *models.OutputMergePullRequestDTO, err error)
	}{
		{
			name: "error GetPullRequestById",
			dto:  &models.InputMergePullRequestDTO{PullRequestId: "PR1"},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().GetPullRequestById(gomock.Any(), "PR1").
					Return(nil, errors.New("db fail"))
			},
			check: func(t *testing.T, out *models.OutputMergePullRequestDTO, err error) {
				assert.Nil(t, out)
				assert.Equal(t, appErrors.ErrServerError, err)
			},
		},
		{
			name: "PR not found",
			dto:  &models.InputMergePullRequestDTO{PullRequestId: "PR1"},
			mockSetup: func(m *mocks.MockRepositoryInterface) {
				m.EXPECT().GetPullRequestById(gomock.Any(), "PR1").
					Return(nil, nil)
			},
			check: func(t *testing.T, out *models.OutputMergePullRequestDTO, err error) {
				assert.Nil(t, out)
				assert.Equal(t, appErrors.ErrResourceNotFound, err)
			},
		},
		{
			name: "PR already merged",
			dto:  &models.InputMergePullRequestDTO{PullRequestId: "PR1"},
			mockSetup: func(m *mocks.MockRepositoryInterface) {

				now := time.Now()

				m.EXPECT().GetPullRequestById(gomock.Any(), "PR1").
					Return(&models.PullRequest{
						SystemId:        "PR1",
						PullRequestName: "Fix bug",
						AuthorSystemId:  "u1",
						Status:          "MERGED",
						AssigneeReviewers: []*models.User{
							{SystemId: "u2"},
							{SystemId: "u3"},
						},
						MergedAt: sql.NullTime{Time: now},
					}, nil)
			},
			check: func(t *testing.T, out *models.OutputMergePullRequestDTO, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "PR1", out.PullRequestID)
				assert.Equal(t, "MERGED", out.Status)
				assert.ElementsMatch(t, []string{"u2", "u3"}, out.AssignedReviewers)
				assert.NotEmpty(t, out.MergedAt)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepositoryInterface(ctrl)
			uc := usecase.NewUseCase(mockRepo)

			tt.mockSetup(mockRepo)

			out, err := uc.MergePullRequest(context.Background(), tt.dto)
			tt.check(t, out, err)
		})
	}
}
