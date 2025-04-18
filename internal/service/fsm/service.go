package fsm

import (
	"context"
	"fmt"
	"kahoot_bsu/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type fSMService struct {
	verificationCodeRepo models.VerificationCodeRepository
	fsmRepo              models.RegistrationFSMRepository
	conn                 *pgxpool.Pool
}

func NewFSMService(
	verificationCodeRepo models.VerificationCodeRepository,
	fsmRepo models.RegistrationFSMRepository,
	conn *pgxpool.Pool,
) *fSMService {
	return &fSMService{
		verificationCodeRepo: verificationCodeRepo,
		fsmRepo:              fsmRepo,
		conn:                 conn,
	}
}

func (f *fSMService) SetState(ctx context.Context, userID int64, code string, state *models.RegistrationFSM) error {
	tx, err := f.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	err = f.verificationCodeRepo.UpdateOrCreate(ctx, userID, code)

	if err != nil {
		//TODO: ADD handle
		return err
	}

	err = f.fsmRepo.UpdateOrCreate(ctx, state)

	if err != nil {
		//TODO: ADD handle
		return err
	}

	return nil
}
