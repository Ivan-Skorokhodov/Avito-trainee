package repository

import (
	"PRmanager/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type RepositoryInterface interface {
	TeamExists(ctx context.Context, teamName string) (bool, error)
	CreateTeam(ctx context.Context, team *models.Team) error
}

func NewRepository() RepositoryInterface {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	rows, err := db.Query(`
	SELECT table_name 
	FROM information_schema.tables 
	WHERE table_schema = 'public' 
	ORDER BY table_name;
`)
	if err != nil {
		log.Fatalf("failed to query tables: %v", err)
	}
	defer rows.Close()

	log.Println("Tables in 'public' schema:")

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			log.Fatalf("failed to scan table name: %v", err)
		}
		log.Printf(" - %s\n", table)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("rows error: %v", err)
	}
	return nil
}
