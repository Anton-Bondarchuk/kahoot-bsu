package game

import (
	"context"
	"kahoot_bsu/internal/domain/models/quiz"
)

// QuizService defines the interface for quiz operations
type QuizService interface {
	GetQuizWithQuestions(ctx context.Context, quizID string) (*quiz.Quiz, error)
}


// // SessionRepository defines the interface for game session operations
// type SessionRepository interface {
// 	Create(ctx context.Context, session *GameSession) error
// 	FindByCode(ctx context.Context, code string) (*GameSession, error)
// 	Update(ctx context.Context, session *GameSession) error
// 	SaveResults(ctx context.Context, session *GameSession) error
// }