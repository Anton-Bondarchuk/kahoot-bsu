package infra

import (
	"context"
	"errors"
	"fmt"
	"kahoot_bsu/internal/domain/question"
	"kahoot_bsu/internal/domain/quiz"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgQuizRepository struct {
	conn *pgxpool.Pool
}

// NewPgQuizRepository creates a new PostgreSQL-based quiz repository
func NewPgQuizRepository(conn *pgxpool.Pool) quiz.Repository {
	return &pgQuizRepository{
		conn: conn,
	}
}

// Update updates an existing quiz with the provided update function
func (r *pgQuizRepository) Update(
	ctx context.Context,
	quizUUID string,
	updateFn func(innerCtx context.Context, quiz *quiz.Quiz) error,
) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	existingQuiz, err := r.getQuizWithTx(ctx, tx, quizUUID)
	if err != nil {
		return err
	}

	err = updateFn(ctx, existingQuiz)
	if err != nil {
		return fmt.Errorf("update function failed: %w", err)
	}

	// Update the quiz
	_, err = tx.Exec(ctx, `
		UPDATE quizzes 
		SET title = $1, updated_at = $2
		WHERE uuid = $3
	`, existingQuiz.Title, time.Now(), quizUUID)
	if err != nil {
		return fmt.Errorf("failed to update quiz: %w", err)
	}

	// Handle questions update
	if err = r.updateQuizQuestions(ctx, tx, existingQuiz); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// UpdateOrCreate updates an existing quiz or creates a new one if it doesn't exist
func (r *pgQuizRepository) UpdateOrCreate(ctx context.Context, quiz *quiz.Quiz) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if quiz exists
	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM quizzes WHERE uuid = $1)", quiz.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if quiz exists: %w", err)
	}

	now := time.Now()
	if exists {
		// Update quiz
		_, err = tx.Exec(ctx, `
			UPDATE quizzes 
			SET title = $1, updated_at = $2
			WHERE uuid = $3
		`, quiz.Title, now, quiz.ID)
		if err != nil {
			return fmt.Errorf("failed to update quiz: %w", err)
		}

		// Remove existing questions and options
		_, err = tx.Exec(ctx, "DELETE FROM questions WHERE quiz_uuid = $1", quiz.ID)
		if err != nil {
			return fmt.Errorf("failed to delete existing questions: %w", err)
		}
	} else {
		// Create new quiz
		_, err = tx.Exec(ctx, `
			INSERT INTO quizzes (uuid, user_uuid, title, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
		`, quiz.ID, quiz.UserID, quiz.Title, now, now)
		if err != nil {
			return fmt.Errorf("failed to create quiz: %w", err)
		}
	}

	// Insert questions and options
	if err = r.updateQuizQuestions(ctx, tx, quiz); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Delete removes a quiz by UUID
func (r *pgQuizRepository) Delete(ctx context.Context, uuid string) error {
	_, err := r.conn.Exec(ctx, "DELETE FROM quizzes WHERE uuid = $1", uuid)
	if err != nil {
		return fmt.Errorf("failed to delete quiz: %w", err)
	}
	return nil
}

// Quiz retrieves a quiz by UUID
func (r *pgQuizRepository) Quiz(ctx context.Context, uuid string) (*quiz.Quiz, error) {
	return r.getQuiz(ctx, uuid)
}

// UserQuizzes retrieves all quizzes belonging to a specific user
func (r *pgQuizRepository) UserQuizzes(ctx context.Context, userID string) ([]*quiz.Quiz, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT uuid, user_uuid, title, created_at, updated_at 
		FROM quizzes
		WHERE user_uuid = $1
		ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user's quizzes: %w", err)
	}
	defer rows.Close()

	var quizzes []*quiz.Quiz
	for rows.Next() {
		q := &quiz.Quiz{}
		if err := rows.Scan(&q.ID, &q.UserID, &q.Title, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan quiz row: %w", err)
		}
		quizzes = append(quizzes, q)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating through quizzes: %w", rows.Err())
	}

	// Load questions for each quiz
	for _, q := range quizzes {
		if err := r.loadQuizQuestions(ctx, q); err != nil {
			return nil, err
		}
	}

	return quizzes, nil
}

// Helper methods

// getQuiz retrieves a quiz by UUID with all its questions and options
func (r *pgQuizRepository) getQuiz(ctx context.Context, uuid string) (*quiz.Quiz, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	quiz, err := r.getQuizWithTx(ctx, tx, uuid)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return quiz, nil
}

// getQuizWithTx retrieves a quiz by UUID within a transaction
func (r *pgQuizRepository) getQuizWithTx(ctx context.Context, tx pgx.Tx, uuid string) (*quiz.Quiz, error) {
	var q quiz.Quiz
	err := tx.QueryRow(ctx, `
		SELECT uuid, user_uuid, title, created_at, updated_at
		FROM quizzes 
		WHERE uuid = $1
	`, uuid).Scan(&q.ID, &q.UserID, &q.Title, &q.CreatedAt, &q.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, quiz.QuizNotFoundError{UUID: uuid}
		}
		return nil, fmt.Errorf("failed to retrieve quiz: %w", err)
	}

	// Load questions with options
	if err := r.loadQuizQuestionsWithTx(ctx, tx, &q); err != nil {
		return nil, err
	}

	return &q, nil
}

// loadQuizQuestions loads questions and options for a quiz
func (r *pgQuizRepository) loadQuizQuestions(ctx context.Context, q *quiz.Quiz) error {
	rows, err := r.conn.Query(ctx, `
		SELECT uuid, quiz_uuid, text, time_limit, points, created_at, updated_at
		FROM questions
		WHERE quiz_uuid = $1
		ORDER BY created_at
	`, q.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch quiz questions: %w", err)
	}
	defer rows.Close()

	return r.scanQuestionsRows(ctx, rows, q)
}

// loadQuizQuestionsWithTx loads questions with transaction
func (r *pgQuizRepository) loadQuizQuestionsWithTx(ctx context.Context, tx pgx.Tx, q *quiz.Quiz) error {
	rows, err := tx.Query(ctx, `
		SELECT uuid, quiz_uuid, text, time_limit, points, created_at, updated_at
		FROM questions
		WHERE quiz_uuid = $1
		ORDER BY created_at
	`, q.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch quiz questions: %w", err)
	}
	defer rows.Close()

	return r.scanQuestionsRows(ctx, rows, q)
}

// scanQuestionsRows scans questions rows and loads options for each question
func (r *pgQuizRepository) scanQuestionsRows(ctx context.Context, rows pgx.Rows, q *quiz.Quiz) error {
	var questions []question.Question
	for rows.Next() {
		var question question.Question
		if err := rows.Scan(
			&question.ID, 
			&question.QuizID, 
			&question.Text, 
			&question.TimeLimit, 
			&question.Points, 
			&question.CreatedAt, 
			&question.UpdatedAt,
		); err != nil {
			return fmt.Errorf("failed to scan question row: %w", err)
		}
		
		// Load options for this question
		if err := r.loadQuestionOptions(ctx, &question); err != nil {
			return err
		}
		
		questions = append(questions, question)
	}

	if rows.Err() != nil {
		return fmt.Errorf("error iterating through questions: %w", rows.Err())
	}

	q.Questions = questions
	return nil
}

// loadQuestionOptions loads options for a question
func (r *pgQuizRepository) loadQuestionOptions(ctx context.Context, question *question.Question) error {
	rows, err := r.conn.Query(ctx, `
		SELECT uuid, question_uuid, text, is_correct
		FROM options
		WHERE question_uuid = $1
	`, question.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch question options: %w", err)
	}
	defer rows.Close()

	var options []quiz.Option
	for rows.Next() {
		var option quiz.Option
		if err := rows.Scan(&option.ID, &option.QuestionID, &option.Text, &option.IsCorrect); err != nil {
			return fmt.Errorf("failed to scan option row: %w", err)
		}
		options = append(options, option)
	}

	if rows.Err() != nil {
		return fmt.Errorf("error iterating through options: %w", rows.Err())
	}

	// TODO: fix type compatibility
	// question.Options = options
	return nil
}

// updateQuizQuestions inserts or updates questions and options for a quiz
func (r *pgQuizRepository) updateQuizQuestions(ctx context.Context, tx pgx.Tx, q *quiz.Quiz) error {
	for i := range q.Questions {
		question := &q.Questions[i]
		question.QuizID = q.ID
		
		// Insert question
		_, err := tx.Exec(ctx, `
			INSERT INTO questions (uuid, quiz_uuid, text, time_limit, points, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, question.ID, question.QuizID, question.Text, question.TimeLimit, question.Points, time.Now(), time.Now())
		if err != nil {
			return fmt.Errorf("failed to insert question: %w", err)
		}

		// Insert options
		for j := range question.Options {
			option := &question.Options[j]
			option.QuestionID = question.ID
			
			_, err := tx.Exec(ctx, `
				INSERT INTO options (uuid, question_uuid, text, is_correct)
				VALUES ($1, $2, $3, $4)
			`, option.QuestionID, option.QuestionID, option.Text, option.IsCorrect)
			if err != nil {
				return fmt.Errorf("failed to insert option: %w", err)
			}
		}
	}
	
	return nil
}