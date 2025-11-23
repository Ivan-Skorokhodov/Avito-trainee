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
		logs.PrintLog(ctx, "[repository] GetTeamByName", appErrors.ErrResourceNotFound.Error())
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
