package infra

import (
	"context"
	"errors"
	"fmt"
	"kahoot_bsu/internal/domain/models/question"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgQuestionRepository struct {
	conn *pgxpool.Pool
}

// NewPgQuestionRepository creates a new PostgreSQL-based question repository
func NewPgQuestionRepository(conn *pgxpool.Pool) question.Repository {
	return &pgQuestionRepository{
		conn: conn,
	}
}

// Create adds a new question to a quiz
func (r *pgQuestionRepository) Create(ctx context.Context, q *question.Question) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Default values
	if q.TimeLimit == 0 {
		q.TimeLimit = 30
	}
	if q.Points == 0 {
		q.Points = 100
	}

	// Insert question
	_, err = tx.Exec(ctx, `
		INSERT INTO questions (uuid, quiz_uuid, text, time_limit, points, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, q.ID, q.QuizID, q.Text, q.TimeLimit, q.Points)
	if err != nil {
		return fmt.Errorf("failed to insert question: %w", err)
	}

	// Insert options
	for i := range q.Options {
		option := &q.Options[i]
		option.QuestionID = q.ID

		_, err = tx.Exec(ctx, `
			INSERT INTO options (uuid, question_uuid, text, is_correct)
			VALUES ($1, $2, $3, $4)
		`, option.ID, option.QuestionID, option.Text, option.IsCorrect)
		if err != nil {
			return fmt.Errorf("failed to insert option: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// Update updates an existing question
func (r *pgQuestionRepository) Update(
	ctx context.Context,
	questionUUID string,
	updateFn func(innerCtx context.Context, question *question.Question) error,
) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get existing question
	existingQuestion, err := r.getQuestionWithTx(ctx, tx, questionUUID)
	if err != nil {
		return err
	}

	// Call update function
	err = updateFn(ctx, existingQuestion)
	if err != nil {
		return fmt.Errorf("update function failed: %w", err)
	}

	// Update the question
	_, err = tx.Exec(ctx, `
		UPDATE questions 
		SET text = $1, time_limit = $2, points = $3, updated_at = $4
		WHERE uuid = $5
	`, existingQuestion.Text, existingQuestion.TimeLimit, existingQuestion.Points, questionUUID)
	if err != nil {
		return fmt.Errorf("failed to update question: %w", err)
	}

	// Update options
	err = r.updateOptionsWithTx(ctx, tx, questionUUID, existingQuestion.Options)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Delete removes a question by UUID
func (r *pgQuestionRepository) Delete(ctx context.Context, uuid string) error {
	// The database should handle cascading deletes for options
	_, err := r.conn.Exec(ctx, "DELETE FROM questions WHERE uuid = $1", uuid)
	if err != nil {
		return fmt.Errorf("failed to delete question: %w", err)
	}
	return nil
}

// Question retrieves a question by UUID
func (r *pgQuestionRepository) Question(ctx context.Context, uuid string) (*question.Question, error) {
	q, err := r.getQuestion(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return q, nil
}

// QuizQuestions retrieves all questions for a specific quiz
func (r *pgQuestionRepository) QuizQuestions(ctx context.Context, quizID string) ([]*question.Question, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT uuid, quiz_uuid, text, time_limit, points, created_at, updated_at
		FROM questions
		WHERE quiz_uuid = $1
		ORDER BY created_at
	`, quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quiz questions: %w", err)
	}
	defer rows.Close()

	return r.scanQuestionsRows(ctx, rows)
}

// UpdateOptions updates the options for a question
func (r *pgQuestionRepository) UpdateOptions(ctx context.Context, questionUUID string, options []question.Option) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// First verify the question exists
	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM questions WHERE uuid = $1)", questionUUID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if question exists: %w", err)
	}
	if !exists {
		return question.QuestionNotFoundError{UUID: questionUUID}
	}

	// Update options
	err = r.updateOptionsWithTx(ctx, tx, questionUUID, options)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Helper methods

// getQuestion retrieves a question by UUID with all its options
func (r *pgQuestionRepository) getQuestion(ctx context.Context, uuid string) (*question.Question, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	q, err := r.getQuestionWithTx(ctx, tx, uuid)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return q, nil
}

// getQuestionWithTx retrieves a question by UUID within a transaction
func (r *pgQuestionRepository) getQuestionWithTx(ctx context.Context, tx pgx.Tx, uuid string) (*question.Question, error) {
	var q question.Question
	err := tx.QueryRow(ctx, `
		SELECT uuid, quiz_uuid, text, time_limit, points, created_at, updated_at
		FROM questions
		WHERE uuid = $1
	`, uuid).Scan(&q.ID, &q.QuizID, &q.Text, &q.TimeLimit, &q.Points)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, question.QuestionNotFoundError{UUID: uuid}
		}
		return nil, fmt.Errorf("failed to retrieve question: %w", err)
	}

	// Load options
	options, err := r.getOptionsWithTx(ctx, tx, uuid)
	if err != nil {
		return nil, err
	}
	q.Options = options

	return &q, nil
}

// scanQuestionsRows scans questions rows and loads options for each question
func (r *pgQuestionRepository) scanQuestionsRows(ctx context.Context, rows pgx.Rows) ([]*question.Question, error) {
	var questions []*question.Question
	for rows.Next() {
		var q question.Question
		if err := rows.Scan(
			&q.ID,
			&q.QuizID,
			&q.Text,
			&q.TimeLimit,
			&q.Points,
		); err != nil {
			return nil, fmt.Errorf("failed to scan question row: %w", err)
		}

		// Load options for this question
		options, err := r.getOptions(ctx, q.ID)
		if err != nil {
			return nil, err
		}
		q.Options = options

		questions = append(questions, &q)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating through questions: %w", rows.Err())
	}

	return questions, nil
}

// getOptions loads options for a question
func (r *pgQuestionRepository) getOptions(ctx context.Context, questionUUID string) ([]question.Option, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT uuid, question_uuid, text, is_correct
		FROM options
		WHERE question_uuid = $1
	`, questionUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch question options: %w", err)
	}
	defer rows.Close()

	return r.scanOptionsRows(rows)
}

// getOptionsWithTx loads options for a question within a transaction
func (r *pgQuestionRepository) getOptionsWithTx(ctx context.Context, tx pgx.Tx, questionUUID string) ([]question.Option, error) {
	rows, err := tx.Query(ctx, `
		SELECT uuid, question_uuid, text, is_correct
		FROM options
		WHERE question_uuid = $1
	`, questionUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch question options: %w", err)
	}
	defer rows.Close()

	return r.scanOptionsRows(rows)
}

// scanOptionsRows scans options rows
func (r *pgQuestionRepository) scanOptionsRows(rows pgx.Rows) ([]question.Option, error) {
	var options []question.Option
	for rows.Next() {
		var opt question.Option
		if err := rows.Scan(
			&opt.ID,
			&opt.QuestionID,
			&opt.Text,
			&opt.IsCorrect,
		); err != nil {
			return nil, fmt.Errorf("failed to scan option row: %w", err)
		}
		options = append(options, opt)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating through options: %w", rows.Err())
	}

	return options, nil
}

// updateOptionsWithTx updates the options for a question within a transaction
func (r *pgQuestionRepository) updateOptionsWithTx(ctx context.Context, tx pgx.Tx, questionUUID string, options []question.Option) error {
	// Delete existing options
	_, err := tx.Exec(ctx, "DELETE FROM options WHERE question_uuid = $1", questionUUID)
	if err != nil {
		return fmt.Errorf("failed to delete existing options: %w", err)
	}

	// Insert new options
	for i := range options {
		option := &options[i]
		option.QuestionID = questionUUID

		_, err = tx.Exec(ctx, `
			INSERT INTO options (uuid, question_uuid, text, is_correct)
			VALUES ($1, $2, $3, $4)
		`, option.ID, option.QuestionID, option.Text, option.IsCorrect)
		if err != nil {
			return fmt.Errorf("failed to insert option: %w", err)
		}
	}

	return nil
}