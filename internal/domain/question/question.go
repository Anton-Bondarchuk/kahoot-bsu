package question

import (
	"context"
	"fmt"
	"time"
)

type Question struct {
	ID       string    `json:"id"`
	QuizID   string    `json:"quiz_id"`
	Text       string    `json:"text"`
	Options    []Option  `json:"options"`
	TimeLimit  int       `json:"time_limit"`
	Points     int       `json:"points"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Option struct {
	ID         string `json:"id"`
	QuestionID string `json:"question_id"`
	Text         string `json:"text"`
	IsCorrect    bool   `json:"is_correct"`
}

type QuestionNotFoundError struct {
	UUID string
}

func (e QuestionNotFoundError) Error() string {
	return fmt.Sprintf("question not found: %s", e.UUID)
}

type Repository interface {
	Create(ctx context.Context, question *Question) error
	Update(ctx context.Context, questionUUID string, updateFn func(innerCtx context.Context, question *Question) error) error
	Delete(ctx context.Context, uuid string) error
	Question(ctx context.Context, uuid string) (*Question, error)
	QuizQuestions(ctx context.Context, quizUUID string) ([]*Question, error)
	UpdateOptions(ctx context.Context, questionUUID string, options []Option) error
}