package handlers

import (
	"context"
	"errors"
	"kahoot_bsu/internal/domain/question"
	"kahoot_bsu/internal/domain/quiz"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handlers contains the HTTP handlers for the API
type Handlers struct {
	quizRepo     quiz.Repository
	questionRepo question.Repository
}

// NewHandlers creates a new Handlers instance
func NewHandlers(quizRepo quiz.Repository, questionRepo question.Repository) *Handlers {
	return &Handlers{
		quizRepo:     quizRepo,
		questionRepo: questionRepo,
	}
}


// GetQuizzes handles GET /api/quizzes
func (h *Handlers) GetQuizzes(c *gin.Context) {
	// ctx := c.Request.Context()
	
	// quizzes, err := h.quizRepo.Quiz(ctx)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quizzes"})
	// 	return
	// }
	
	// c.JSON(http.StatusOK, quizzes)
}

// CreateQuiz handles POST /api/quizzes
func (h *Handlers) CreateQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	
	var quizData quiz.Quiz
	if err := c.ShouldBindJSON(&quizData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Generate a new UUID if not provided
	if quizData.ID == "" {
		quizData.ID = uuid.NewString()
	}
	
	// Set user UUID from auth context
	userUUID := c.GetString("userUUID")
	if userUUID == "" {
		userUUID = "anonymous" // Fallback
	}
	quizData.UserID = userUUID
	
	// Set timestamps
	now := time.Now()
	quizData.CreatedAt = now
	quizData.UpdatedAt = now
	
	if err := h.quizRepo.UpdateOrCreate(ctx, &quizData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quiz"})
		return
	}
	
	c.JSON(http.StatusCreated, quizData)
}

// GetQuiz handles GET /api/quizzes/:id
func (h *Handlers) GetQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	
	quizUUID := c.Param("id")
	if quizUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing quiz ID"})
		return
	}
	
	quizData, err := h.quizRepo.Quiz(ctx, quizUUID)
	if err != nil {
		var quizNotFoundErr quiz.QuizNotFoundError
		if errors.As(err, &quizNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quiz"})
		}
		return
	}
	
	c.JSON(http.StatusOK, quizData)
}

// UpdateQuiz handles PUT /api/quizzes/:id
func (h *Handlers) UpdateQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	
	quizUUID := c.Param("id")
	if quizUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing quiz ID"})
		return
	}
	
	var updatedQuiz quiz.Quiz
	if err := c.ShouldBindJSON(&updatedQuiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Ensure we use the UUID from the URL
	updatedQuiz.ID = quizUUID
	updatedQuiz.UpdatedAt = time.Now()
	
	err := h.quizRepo.Update(ctx, quizUUID, func(innerCtx context.Context, q *quiz.Quiz) error {
		q.Title = updatedQuiz.Title
		// Update other fields as needed but preserve metadata
		return nil
	})
	
	if err != nil {
		var quizNotFoundErr quiz.QuizNotFoundError
		if errors.As(err, &quizNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quiz"})
		}
		return
	}
	
	// Get the updated quiz to return
	quizData, err := h.quizRepo.Quiz(ctx, quizUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated quiz"})
		return
	}
	
	c.JSON(http.StatusOK, quizData)
}

// DeleteQuiz handles DELETE /api/quizzes/:id
func (h *Handlers) DeleteQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	
	quizUUID := c.Param("id")
	if quizUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing quiz ID"})
		return
	}
	
	err := h.quizRepo.Delete(ctx, quizUUID)
	if err != nil {
		var quizNotFoundErr quiz.QuizNotFoundError
		if errors.As(err, &quizNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete quiz"})
		}
		return
	}
	
	c.Status(http.StatusNoContent)
}

// GetQuizQuestions handles GET /api/quizzes/:id/questions
func (h *Handlers) GetQuizQuestions(c *gin.Context) {
	ctx := c.Request.Context()
	
	quizUUID := c.Param("id")
	if quizUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing quiz ID"})
		return
	}
	
	// Verify quiz exists
	_, err := h.quizRepo.Quiz(ctx, quizUUID)
	if err != nil {
		var quizNotFoundErr quiz.QuizNotFoundError
		if errors.As(err, &quizNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify quiz"})
		}
		return
	}
	
	questions, err := h.questionRepo.QuizQuestions(ctx, quizUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}
	
	c.JSON(http.StatusOK, questions)
}

// AddQuizQuestion handles POST /api/quizzes/:id/questions
func (h *Handlers) AddQuizQuestion(c *gin.Context) {
	ctx := c.Request.Context()
	
	quizUUID := c.Param("id")
	if quizUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing quiz ID"})
		return
	}
	
	// Verify quiz exists
	_, err := h.quizRepo.Quiz(ctx, quizUUID)
	if err != nil {
		var quizNotFoundErr quiz.QuizNotFoundError
		if errors.As(err, &quizNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify quiz"})
		}
		return
	}
	
	var questionData question.Question
	if err := c.ShouldBindJSON(&questionData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Generate a new UUID for the question
	questionData.ID = uuid.NewString()
	questionData.QuizID = quizUUID
	
	// Set timestamps
	now := time.Now()
	questionData.CreatedAt = now
	questionData.UpdatedAt = now
	
	// Generate UUIDs for options if needed
	for i := range questionData.Options {
		if questionData.Options[i].ID == "" {
			questionData.Options[i].ID = uuid.NewString()
		}
		questionData.Options[i].QuestionID = questionData.ID
	}
	
	if err := h.questionRepo.Create(ctx, &questionData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}
	
	c.JSON(http.StatusCreated, questionData)
}

// UpdateQuestion handles PUT /api/questions/:question_id
func (h *Handlers) UpdateQuestion(c *gin.Context) {
	ctx := c.Request.Context()
	
	questionUUID := c.Param("question_id")
	if questionUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing question ID"})
		return
	}
	
	var updatedQuestion question.Question
	if err := c.ShouldBindJSON(&updatedQuestion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Ensure we use the UUID from the URL
	updatedQuestion.ID = questionUUID
	
	err := h.questionRepo.Update(ctx, questionUUID, func(innerCtx context.Context, q *question.Question) error {
		q.Text = updatedQuestion.Text
		q.TimeLimit = updatedQuestion.TimeLimit
		q.Points = updatedQuestion.Points
		q.Options = updatedQuestion.Options
		q.UpdatedAt = time.Now()
		
		// Ensure all options have UUIDs and the correct question UUID
		for i := range q.Options {
			if q.Options[i].ID == "" {
				q.Options[i].ID = uuid.NewString()
			}
			q.Options[i].QuestionID = questionUUID
		}
		
		return nil
	})
	
	if err != nil {
		var questionNotFoundErr question.QuestionNotFoundError
		if errors.As(err, &questionNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question"})
		}
		return
	}
	
	// Get the updated question to return
	updatedData, err := h.questionRepo.Question(ctx, questionUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated question"})
		return
	}
	
	c.JSON(http.StatusOK, updatedData)
}

// DeleteQuestion handles DELETE /api/questions/:question_id
func (h *Handlers) DeleteQuestion(c *gin.Context) {
	ctx := c.Request.Context()
	
	questionUUID := c.Param("question_id")
	if questionUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing question ID"})
		return
	}
	
	err := h.questionRepo.Delete(ctx, questionUUID)
	if err != nil {
		var questionNotFoundErr question.QuestionNotFoundError
		if errors.As(err, &questionNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete question"})
		}
		return
	}
	
	c.Status(http.StatusNoContent)
}