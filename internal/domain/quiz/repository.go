package quiz

import (
	"context"
	"fmt"
)

type QuizNotFoundError struct {
	UUID string
}

func (e QuizNotFoundError) Error() string {
	return fmt.Sprintf("bot not found: %s", e.UUID)
}

type Repository interface {
	Update(
		ctx context.Context,
		quizUUID string,
		updateFn func(innerCtx context.Context, bot *Quiz) error,
	) error
	UpdateOrCreate(ctx context.Context, bot *Quiz) error
	Delete(ctx context.Context, uuid string) error
	Quiz(ctx context.Context, uuid string) (*Quiz, error)
	UserQuizzes(ctx context.Context, userUUID string) ([]*Quiz, error)
}