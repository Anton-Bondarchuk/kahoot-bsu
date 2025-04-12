package player

import (
	"time"

	"github.com/google/uuid"
)

// Player represents a user playing a quiz
type Player struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	GameID   uuid.UUID `json:"game_id"`
	Active   bool      `json:"active"`
	JoinedAt time.Time `json:"joined_at"`
}

// Game represents a quiz game session
type Game struct {
	ID           uuid.UUID    `json:"id"`
	QuizID       uuid.UUID    `json:"quiz_id"`
	HostID       uuid.UUID    `json:"host_id"`
	Code         string       `json:"code"`
	Status       GameStatus   `json:"status"`
	CurrentIndex int          `json:"current_index"`
	StartedAt    *time.Time   `json:"started_at"`
	EndedAt      *time.Time   `json:"ended_at"`
	Players      []Player     `json:"players"`
	Questions    []Question   `json:"questions"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// GameStatus represents the status of a game
type GameStatus string

const (
	GameStatusWaiting  GameStatus = "waiting"
	GameStatusStarted  GameStatus = "started"
	GameStatusQuestion GameStatus = "question"
	GameStatusReview   GameStatus = "review"
	GameStatusEnded    GameStatus = "ended"
)

// Quiz represents a collection of questions
type Quiz struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatorID   uuid.UUID  `json:"creator_id"`
	Questions   []Question `json:"questions"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Question represents a quiz question
type Question struct {
	ID         uuid.UUID `json:"id"`
	QuizID     uuid.UUID `json:"quiz_id"`
	Text       string    `json:"text"`
	TimeLimit  int       `json:"time_limit"` // in seconds
	Points     int       `json:"points"`
	Index      int       `json:"index"`
	Options    []Option  `json:"options"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Option represents an answer option for a question
type Option struct {
	ID         uuid.UUID `json:"id"`
	QuestionID uuid.UUID `json:"question_id"`
	Text       string    `json:"text"`
	IsCorrect  bool      `json:"is_correct,omitempty"` // Only visible to admins
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Answer represents a player's answer to a question
type Answer struct {
	ID         uuid.UUID `json:"id"`
	PlayerID   uuid.UUID `json:"player_id"`
	GameID     uuid.UUID `json:"game_id"`
	QuestionID uuid.UUID `json:"question_id"`
	OptionID   uuid.UUID `json:"option_id"`
	Timestamp  time.Time `json:"timestamp"`
	TimeToAnswer int      `json:"time_to_answer"` // in milliseconds
}

// PlayerGameSummary represents a player's game summary
type PlayerGameSummary struct {
	PlayerID  uuid.UUID `json:"player_id"`
	GameID    uuid.UUID `json:"game_id"`
	Score     int       `json:"score"`
	Position  int       `json:"position"`
	Correct   int       `json:"correct"`
	Incorrect int       `json:"incorrect"`
}

// Event types for RabbitMQ events
const (
	EventPlayerJoined       = "player.joined"
	EventPlayerLeft         = "player.left"
	EventGameCreated        = "game.created"
	EventGameStarted        = "game.started"
	EventGameEnded          = "game.ended"
	EventQuestionStarted    = "question.started"
	EventQuestionEnded      = "question.ended"
	EventAnswerSubmitted    = "answer.submitted"
	EventLeaderboardUpdated = "leaderboard.updated"
)

// Message represents a message to be published to RabbitMQ
type Message struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

// Reponse structs for client communications
type GameStateResponse struct {
	Status       GameStatus  `json:"status"`
	Question     *Question   `json:"question,omitempty"`
	TimeLeft     int         `json:"time_left,omitempty"`
	Players      []Player    `json:"players,omitempty"`
	Leaderboard  interface{} `json:"leaderboard,omitempty"`
	CurrentIndex int         `json:"current_index,omitempty"`
	TotalCount   int         `json:"total_count,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}