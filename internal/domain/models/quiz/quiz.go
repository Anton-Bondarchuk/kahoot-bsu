package quiz

import (
	"kahoot_bsu/internal/domain/models/question"
	"time"
)

type Quiz struct {
	ID        string       `json:"id"`
	UserID	  string
	Title     string     `json:"title" binding:"required"`
	IsPublic    bool     `json:"is_public"`
	CreatedBy string     `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Questions []question.Question `json:"questions,omitempty"`
}