package infra

import (
	"context"
	"fmt"
	"kahoot_bsu/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgFsmRegistrationRepository struct {
	conn *pgxpool.Pool
}

func NewPgFsmRegistrationRepository(conn *pgxpool.Pool) models.RegistrationFSMRepository {
	return &pgFsmRegistrationRepository{
		conn: conn,
	}
}

func (r *pgFsmRegistrationRepository) UpdateOrCreate(ctx context.Context, fsm *models.RegistrationFSM) error {
	var exists bool
	err := r.conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM registration_fsm WHERE user_id = $1)", fsm.UserID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if FSM exists: %w", err)
	}
	
	if exists {
		_, err := r.conn.Exec(ctx, `
			UPDATE registration_fsm
			SET 
				wait_login = $1, 
				wait_otp = $2,
				is_registered = $3,
			WHERE user_id = $4
		`, fsm.WaitLogin, fsm.WaitOTP, fsm.IsRegistered, fsm.UserID)
		
		if err != nil {
			return fmt.Errorf("failed to update FSM: %w", err)
		}

		return nil
	} 

	fsmID := uuid.NewString()

	_, err = r.conn.Exec(ctx, `
		INSERT INTO registration_fsm (id, user_id, wait_login, wait_otp, is_registered)
		VALUES ($1, $2, $3, $4, $5)
	`, fsmID, fsm.UserID, fsm.WaitLogin, fsm.WaitOTP, fsm.IsRegistered)

	if err != nil {
		return fmt.Errorf("failed to create FSM: %w", err)
	}

	return nil
}