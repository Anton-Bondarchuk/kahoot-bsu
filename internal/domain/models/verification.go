package models

import (
	"context"
	"time"
)

type VerificationCode struct {
	ID        string
	UserID    int64
	Code      string
	ExprireAt time.Time
	CreatedAt time.Time
}

type VerificationCodeRepository interface {
	UpdateOrCreate(ctx context.Context, userId int64, code string) error
	Delete(ctx context.Context, deleteFn func(innerCtx context.Context)) error
	DeleteByUserId(ctx context.Context, userId int64) error
}
