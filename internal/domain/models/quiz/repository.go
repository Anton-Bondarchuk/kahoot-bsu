package quiz

import (
	"context"
	"fmt"
)

type QuizNotFoundError struct {
	ID string
}

func (e QuizNotFoundError) Error() string {
	return fmt.Sprintf("bot not found: %s", e.ID)
}

type Repository interface {
	Update(
		ctx context.Context,
		quizID string,
		updateFn func(innerCtx context.Context, quiz *Quiz) error,
	) error
	UpdateOrCreate(ctx context.Context, quiz *Quiz) error
	Delete(ctx context.Context, uuid string) error
	Quiz(ctx context.Context, uuid string) (*Quiz, error)
	UserQuizzes(ctx context.Context, userId int64) ([]*Quiz, error)
}