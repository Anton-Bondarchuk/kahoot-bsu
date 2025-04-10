package quiz

import (
	"kahoot_bsu/internal/domain/question"
	"time"
)

// Quiz represents a quiz entity with questions
type Quiz struct {
	ID        string       `json:"id"`
	UserID	  string
	Title     string     `json:"title" binding:"required"`
	CreatedBy string     `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Questions []question.Question `json:"questions,omitempty"`
}

type Option struct {
	ID         uint   `json:"id"`
	QuestionID uint   `json:"question_id"`
	Text       string `json:"text" binding:"required"`
	IsCorrect  bool   `json:"is_correct"`
}
