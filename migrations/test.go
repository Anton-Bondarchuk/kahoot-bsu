package migrations

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)


func createTables(ctx context.Context, db *pgxpool.Pool) error {
	// Create quizzes table
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS quizzes (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			created_by VARCHAR(255),
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create questions table
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS questions (
			id SERIAL PRIMARY KEY,
			quiz_id INTEGER NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
			text TEXT NOT NULL,
			time_limit INTEGER DEFAULT 30,
			points INTEGER DEFAULT 100,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create options table
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS options (
			id SERIAL PRIMARY KEY,
			question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
			text TEXT NOT NULL,
			is_correct BOOLEAN NOT NULL DEFAULT FALSE
		)
	`)
	return err
}