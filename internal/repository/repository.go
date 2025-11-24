package repository

import (
	"PRmanager/internal/models"
	appErrors "PRmanager/pkg/app_errors"
	"PRmanager/pkg/logs"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type RepositoryInterface interface {
	TeamExists(ctx context.Context, teamName string) (bool, error)
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeamByName(ctx context.Context, teamName string) (*models.Team, error)
	SetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
	GetUserBySystemId(ctx context.Context, systemId string) (*models.User, error)
	GetListReviewsByUserId(ctx context.Context, userId int) ([]*models.PullRequest, error)
	PullRequestExists(ctx context.Context, prSystemID string) (bool, error)
	GetTeamMembers(ctx context.Context, teamId int) ([]*models.User, error)
	CreatePullRequestAndReview(ctx context.Context, pr *models.PullRequest, reviews []*models.User) error
	GetPullRequestById(ctx context.Context, prSystemId string) (*models.PullRequest, error)
	SetMergedStatusPullRequest(ctx context.Context, prId int) (sql.NullTime, error)
	ReplaceReviewers(ctx context.Context, prId int, oldReviewerId int, newReviewerId int) error
}

type Database struct {
	conn *sql.DB
}

func NewDatabase() *Database {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	return &Database{conn: conn}
}

func (db *Database) TeamExists(ctx context.Context, teamName string) (bool, error) {
	const query = `
        SELECT EXISTS(
            SELECT 1 FROM teams WHERE team_name = $1
        );
    `

	var exists bool
	err := db.conn.QueryRowContext(ctx, query, teamName).Scan(&exists)
	if err != nil {
		logs.PrintLog(ctx, "[repository] TeamExists", err.Error())
		return false, err
	}

	return exists, nil
}

func (db *Database) CreateTeam(ctx context.Context, team *models.Team) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		logs.PrintLog(ctx, "[repository] CreateTeam", err.Error())
		return err
	}

	const insertTeam = `
        INSERT INTO teams (team_name)
        VALUES ($1)
        RETURNING team_id;
    `
	if err := tx.QueryRowContext(ctx, insertTeam, team.TeamName).Scan(&team.TeamId); err != nil {
		tx.Rollback()
		logs.PrintLog(ctx, "[repository] CreateTeam", err.Error())
		return err
	}

	const insertUser = `
        INSERT INTO users (system_id, user_name, team_id, is_active)
        VALUES ($1, $2, $3, $4)
        RETURNING user_id;
    `

	for _, member := range team.TeamMembers {
		var newUserID int
		err := tx.QueryRowContext(
			ctx,
			insertUser,
			member.SystemId,
			member.UserName,
			team.TeamId,
			member.IsActive,
		).Scan(&newUserID)

		if err != nil {
			tx.Rollback()
			logs.PrintLog(ctx, "[repository] CreateTeam", err.Error())
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		logs.PrintLog(ctx, "[repository] CreateTeam", err.Error())
		return err
	}

	return nil
}

func (db *Database) GetTeamByName(ctx context.Context, teamName string) (*models.Team, error) {
	const selectTeam = `
        SELECT team_id, team_name
        FROM teams
        WHERE team_name = $1;
    `

	var team models.Team

	err := db.conn.QueryRowContext(ctx, selectTeam, teamName).
		Scan(&team.TeamId, &team.TeamName)

	if errors.Is(err, sql.ErrNoRows) {
		logs.PrintLog(ctx, "[repository] GetTeamByName", err.Error())
		return nil, nil
	}

	if err != nil {
		logs.PrintLog(ctx, "[repository] GetTeamByName", err.Error())
		return nil, err
	}

	const selectMembers = `
        SELECT user_id, system_id, user_name, team_id, is_active
        FROM users
        WHERE team_id = $1;
    `

	rows, err := db.conn.QueryContext(ctx, selectMembers, team.TeamId)
	if err != nil {
		logs.PrintLog(ctx, "[repository] GetTeamByName", err.Error())
		return nil, err
	}
	defer rows.Close()

	team.TeamMembers = make([]*models.User, 0)
	for rows.Next() {
		member := &models.User{}

		err := rows.Scan(
			&member.UserId,
			&member.SystemId,
			&member.UserName,
			&member.TeamId,
			&member.IsActive,
		)

		if err != nil {
			logs.PrintLog(ctx, "[repository] GetTeamByName", err.Error())
			return nil, err
		}

		team.TeamMembers = append(team.TeamMembers, member)
	}

	return &team, nil
}

func (db *Database) SetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	const query = `
        UPDATE users
        SET is_active = $2
        WHERE system_id = $1
        RETURNING 
            users.user_id,
            users.system_id,
            users.user_name,
            users.team_id,
            users.is_active;
    `

	var user models.User
	err := db.conn.
		QueryRowContext(ctx, query, userID, isActive).
		Scan(&user.UserId, &user.SystemId, &user.UserName, &user.TeamId, &user.IsActive)

	if errors.Is(err, sql.ErrNoRows) {
		logs.PrintLog(ctx, "[repository] SetIsActive", err.Error())
		return nil, nil
	}

	if err != nil {
		logs.PrintLog(ctx, "[repository] SetIsActive", err.Error())
		return nil, appErrors.ErrServerError
	}

	const teamQuery = `
            SELECT team_name 
            FROM teams 
            WHERE team_id = $1;
        `
	err = db.conn.QueryRowContext(ctx, teamQuery, user.TeamId).Scan(&user.TeamName)

	if err != nil {
		logs.PrintLog(ctx, "[repository] SetIsActive", err.Error())
		return nil, err
	}

	return &user, nil
}

func (db *Database) GetUserBySystemId(ctx context.Context, systemId string) (*models.User, error) {
	const userQuery = `
        SELECT 
            user_id,
			team_id
        FROM users
        WHERE system_id = $1;
    `

	var user models.User
	user.SystemId = systemId
	err := db.conn.QueryRowContext(ctx, userQuery, systemId).Scan(&user.UserId, &user.TeamId)

	if errors.Is(err, sql.ErrNoRows) {
		logs.PrintLog(ctx, "[repository] GetUserBySystemId", err.Error())
		return nil, nil
	}

	if err != nil {
		logs.PrintLog(ctx, "[repository] GetUserBySystemId", err.Error())
		return nil, err
	}
	return &user, nil
}

func (db *Database) GetListReviewsByUserId(ctx context.Context, userId int) ([]*models.PullRequest, error) {
	const prQuery = `
        SELECT
            pr.system_id,
            pr.pull_request_name,
            au.system_id,
            pr.status
        FROM pull_request_reviewers AS r
        JOIN pull_requests AS pr ON pr.pull_request_id = r.pull_request_id
        JOIN users AS au ON au.user_id = pr.author_id
        WHERE r.user_id = $1;
    `

	rows, err := db.conn.QueryContext(ctx, prQuery, userId)
	if err != nil {
		logs.PrintLog(ctx, "[repository] GetListReviewsByUserId", err.Error())
		return nil, err
	}
	defer rows.Close()

	reviews := make([]*models.PullRequest, 0)
	for rows.Next() {
		pr := &models.PullRequest{}

		err := rows.Scan(
			&pr.SystemId,
			&pr.PullRequestName,
			&pr.AuthorSystemId,
			&pr.Status,
		)
		if err != nil {
			logs.PrintLog(ctx, "[repository] GetListReviewsByUserId", err.Error())
			return nil, err
		}

		reviews = append(reviews, pr)
	}
	return reviews, nil
}

func (db *Database) PullRequestExists(ctx context.Context, prSystemID string) (bool, error) {
	const query = `
        SELECT EXISTS (
            SELECT 1 
            FROM pull_requests
            WHERE system_id = $1
        );
    `

	var exists bool
	err := db.conn.QueryRowContext(ctx, query, prSystemID).Scan(&exists)
	if err != nil {
		logs.PrintLog(ctx, "[repository] PullRequestExists", err.Error())
		return false, err
	}

	return exists, nil
}

func (db *Database) GetTeamMembers(ctx context.Context, teamId int) ([]*models.User, error) {
	const query = `
        SELECT 
            user_id,
            system_id,
            user_name,
            is_active
        FROM users
        WHERE team_id = $1;
    `

	rows, err := db.conn.QueryContext(ctx, query, teamId)
	if err != nil {
		logs.PrintLog(ctx, "[repository] GetTeamMembers", err.Error())
		return nil, err
	}
	defer rows.Close()

	members := make([]*models.User, 0)

	for rows.Next() {
		m := &models.User{}

		err := rows.Scan(
			&m.UserId,
			&m.SystemId,
			&m.UserName,
			&m.IsActive,
		)
		if err != nil {
			logs.PrintLog(ctx, "[repository] GetTeamMembers", err.Error())
			return nil, err
		}

		members = append(members, m)
	}

	return members, nil
}

func (db *Database) CreatePullRequestAndReview(ctx context.Context, pr *models.PullRequest, reviewers []*models.User) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		logs.PrintLog(ctx, "[repository] CreatePullRequestAndReview", err.Error())
		return err
	}

	const insertPR = `
        INSERT INTO pull_requests (system_id, pull_request_name, author_id, status)
        VALUES ($1, $2, $3, $4)
        RETURNING pull_request_id;
    `

	err = tx.QueryRowContext(
		ctx,
		insertPR,
		pr.SystemId,
		pr.PullRequestName,
		pr.AuthorId,
		pr.Status,
	).Scan(&pr.PullRequestId)

	if err != nil {
		tx.Rollback()
		logs.PrintLog(ctx, "[repository] CreatePullRequestAndReview", err.Error())
		return err
	}

	const insertReviewer = `
        INSERT INTO pull_request_reviewers (pull_request_id, user_id)
        VALUES ($1, $2)
        ON CONFLICT DO NOTHING;
    `

	for _, r := range reviewers {
		_, err := tx.ExecContext(ctx, insertReviewer, pr.PullRequestId, r.UserId)
		if err != nil {
			tx.Rollback()
			logs.PrintLog(ctx, "[repository] CreatePullRequestAndReview", err.Error())
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logs.PrintLog(ctx, "[repository] CreatePullRequestAndReview", err.Error())
		return err
	}

	return nil
}

func (db *Database) GetPullRequestById(ctx context.Context, prSystemId string) (*models.PullRequest, error) {
	const prQuery = `
        SELECT 
            pr.pull_request_id,
            pr.system_id,
            pr.pull_request_name,
            pr.author_id,
            au.system_id AS author_system_id,
            pr.status,
            pr.created_at,
            pr.merged_at
        FROM pull_requests AS pr
        JOIN users AS au ON au.user_id = pr.author_id
        WHERE pr.system_id = $1;
    `

	pr := &models.PullRequest{}
	err := db.conn.QueryRowContext(ctx, prQuery, prSystemId).Scan(
		&pr.PullRequestId,
		&pr.SystemId,
		&pr.PullRequestName,
		&pr.AuthorId,
		&pr.AuthorSystemId,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		logs.PrintLog(ctx, "[repository] GetPullRequestById", err.Error())
		return nil, err
	}

	const reviewersQuery = `
        SELECT 
            u.user_id,
            u.system_id,
            u.user_name
        FROM pull_request_reviewers AS r
        JOIN users AS u ON u.user_id = r.user_id
        WHERE r.pull_request_id = $1;
    `

	rows, err := db.conn.QueryContext(ctx, reviewersQuery, pr.PullRequestId)
	if err != nil {
		logs.PrintLog(ctx, "[repository] GetPullRequestById", err.Error())
		return nil, err
	}
	defer rows.Close()

	pr.AssigneeReviewers = make([]*models.User, 0)

	for rows.Next() {
		u := &models.User{}
		err := rows.Scan(
			&u.UserId,
			&u.SystemId,
			&u.UserName,
		)

		if err != nil {
			logs.PrintLog(ctx, "[repository] GetPullRequestById", err.Error())
			return nil, err
		}

		pr.AssigneeReviewers = append(pr.AssigneeReviewers, u)
	}

	return pr, nil
}

func (db *Database) SetMergedStatusPullRequest(ctx context.Context, prId int) (sql.NullTime, error) {
	const query = `
        UPDATE pull_requests
        SET 
            status = 'MERGED',
            merged_at = NOW()
        WHERE pull_request_id = $1
        RETURNING merged_at;
    `

	var mergedAt sql.NullTime

	err := db.conn.QueryRowContext(ctx, query, prId).Scan(&mergedAt)
	if err != nil {
		logs.PrintLog(ctx, "[repository] SetMergedStatusPullRequest", err.Error())
		return sql.NullTime{}, err
	}

	return mergedAt, nil
}

func (db *Database) ReplaceReviewers(ctx context.Context, prId int, oldReviewerId int, newReviewerId int) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		logs.PrintLog(ctx, "[repository] ReplaceReviewers", err.Error())
		return err
	}

	const deleteQuery = `
        DELETE FROM pull_request_reviewers
        WHERE pull_request_id = $1 AND user_id = $2;
    `

	result, err := tx.ExecContext(ctx, deleteQuery, prId, oldReviewerId)
	if err != nil {
		tx.Rollback()
		logs.PrintLog(ctx, "[repository] ReplaceReviewers", err.Error())
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		tx.Rollback()
		logs.PrintLog(ctx, "[repository] ReplaceReviewers", err.Error())
		return err
	}

	const insertQuery = `
        INSERT INTO pull_request_reviewers (pull_request_id, user_id)
        VALUES ($1, $2)
        ON CONFLICT DO NOTHING;
    `

	result, err = tx.ExecContext(ctx, insertQuery, prId, newReviewerId)
	if err != nil {
		tx.Rollback()
		logs.PrintLog(ctx, "[repository] ReplaceReviewers", err.Error())
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		tx.Rollback()
		logs.PrintLog(ctx, "[repository] ReplaceReviewers", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		logs.PrintLog(ctx, "[repository] ReplaceReviewers", err.Error())
		return err
	}

	logs.PrintLog(ctx, "[repository] ReplaceReviewers", fmt.Sprint("success replace: %+v -> %+v", oldReviewerId, newReviewerId))
	return nil
}
