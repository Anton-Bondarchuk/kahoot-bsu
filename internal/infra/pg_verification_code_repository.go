package infra

import (
	"context"
	"fmt"
	"kahoot_bsu/internal/domain/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgVerificationCodeRepository struct {
	conn *pgxpool.Pool
}

func NewPgVerificationCodeRepository(conn *pgxpool.Pool) models.VerificationCodeRepository {
	return &pgVerificationCodeRepository{
		conn: conn,
	}
}

// Create implements models.VerificationCodeRepository.
func (r *pgVerificationCodeRepository) UpdateOrCreate(ctx context.Context, userId int64, code string) error {
	// TODO: declare the start errors
	// specify the operation
	// add logger: chekc sql-injection code
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)

	}
	defer tx.Rollback(ctx)

	var exist bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM quizzes WHERE user_id = $1)", userId).Scan(&exist)

	expiresAat := time.Now().Local().Add(time.Minute * 30)
	if exist {
		_, err = tx.Exec(ctx, `
		UPDATE verification_codes
		SET code = $1, expires_at = $2
		WHERE user_id = $3
		`, code, expiresAat, userId)

		if err != nil {
			return fmt.Errorf("failed to verification code quiz: %w", err)
		}

		return tx.Commit(ctx)
	}

	id := uuid.NewString()
	now := time.Now()
	_, err = tx.Exec(ctx, `
		INSERT INTO verification_codes (id, user_id, code, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, id, userId, code, expiresAat, now)
	if err != nil {
		return fmt.Errorf("failed to create verification code: %w", err)
	}

	return tx.Commit(ctx)
}

// Delete implements models.VerificationCodeRepository.
func (r *pgVerificationCodeRepository) Delete(ctx context.Context, deleteFn func(innerCtx context.Context)) error {
	panic("pgVerificationCodeRepository")
}

// DeleteByUserId implements models.VerificationCodeRepository.
func (r *pgVerificationCodeRepository) DeleteByUserId(ctx context.Context, userId int64) error {
	_, err := r.conn.Exec(ctx, "DELETE FROM verification_codes WHERE user_id = $1", userId)
	if err != nil {
		return fmt.Errorf("failed to delete verification code: %w", err)
	}
	return nil
}
