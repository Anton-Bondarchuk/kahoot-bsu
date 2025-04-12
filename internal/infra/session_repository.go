package infra

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"kahoot_bsu/internal/kahoot/game"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// var (
// 	GameSessionNotFoundErr = errors.New("game session not found")
// )

// type PgSessionRepository struct {
// 	db *pgxpool.Pool
// }

// func NewPgSessionRepository(db *pgxpool.Pool) *PgSessionRepository {
// 	return &PgSessionRepository{
// 		db: db,
// 	}
// }

// // Create stores a new game session in the database
// func (r *PgSessionRepository) Create(ctx context.Context, session *game.GameSession) error {
// 	query := `
// 		INSERT INTO game_sessions (
// 			id, quiz_id, host_id, join_code, status_flags, 
// 			current_question_index, started_at, created_at
// 		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
// 	`

// 	_, err := r.db.Exec(ctx, query,
// 		session.ID,
// 		session.QuizID,
// 		session.HostID,
// 		session.Code,
// 		session.Status,
// 		session.CurrentQuestionIndex,
// 		nil, 
// 		session.CreatedAt,
// 	)

// 	return err
// }

// func (r *PgSessionRepository) FindByCode(ctx context.Context, code string) (*game.GameSession, error) {
// 	query := `
// 		SELECT 
// 			id, quiz_id, host_id, join_code, status, 
// 			current_question_index, started_at, ended_at, created_at
// 		FROM game_sessions 
// 		WHERE join_code = $1
// 	`

// 	row := r.db.QueryRow(ctx, query, code)

// 	var session game.GameSession
// 	var startedAt, endedAt sql.NullTime

// 	err := row.Scan(
// 		&session.ID,
// 		&session.QuizID,
// 		&session.HostID,
// 		&session.Code,
// 		&session.Status,
// 		&session.CurrentQuestionIndex,
// 		&startedAt,
// 		&endedAt,
// 		&session.CreatedAt,
// 	)

// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return nil, GameSessionNotFoundErr
// 		}
// 		return nil, err
// 	}

// 	if startedAt.Valid {
// 		session.StartedAt = startedAt.Time
// 	}

// 	if endedAt.Valid {
// 		session.EndedAt = endedAt.Time
// 	}

// 	// Initialize players map
// 	session.Players = make(map[string]*game.Player)

// 	// Fetch the quiz for this session
// 	// This would typically be done by a quiz repository or service
// 	// Here we just create a placeholder for the example
// 	session.Quiz = &game.Quiz{
// 		ID:    session.QuizID,
// 		Title: "Sample Quiz", // This would be fetched from the database
// 	}

// 	return &session, nil
// }

// func (r *PgSessionRepository) Update(ctx context.Context, session *game.GameSession) error {
// 	query := `
// 		UPDATE game_sessions 
// 		SET 
// 			status = $1, 
// 			current_question_index = $2, 
// 			started_at = $3, 
// 			ended_at = $4
// 		WHERE id = $5
// 	`

// 	startedAt := session.StartedAt
// 	endedAt := session.EndedAt

// 	_, err := r.db.Exec(ctx, query,
// 		session.Status,
// 		session.CurrentQuestionIndex,
// 		startedAt,
// 		endedAt,
// 		session.ID,
// 	)

// 	return err
// }

// func (r *PgSessionRepository) SaveResults(ctx context.Context, session *game.GameSession) error {
// 	tx, err := r.db.Begin(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback(ctx)

// 	_, err = tx.Exec(ctx, `
// 		UPDATE game_sessions 
// 		SET 
// 			status = $1, 
// 			ended_at = $2
// 		WHERE id = $3
// 	`, "finished", time.Now(), session.ID)

// 	if err != nil {
// 		return err
// 	}

// 	for playerID, player := range session.Players {
// 		var participantID string

// 		err = tx.QueryRow(ctx, `
// 			INSERT INTO participants (
// 				id, session_id, user_id, nickname, score, joined_at
// 			) VALUES (
// 				$1, $2, $3, $4, $5, $6
// 			)
// 			ON CONFLICT (id) DO UPDATE 
// 			SET score = $5
// 			RETURNING id
// 		`, playerID, session.ID, player.UserID, player.Nickname, player.Score, time.Now()).Scan(&participantID)

// 		if err != nil {
// 			return err
// 		}

// 		for questionID, answer := range player.Answers {
// 			_, err = tx.Exec(ctx, `
// 				INSERT INTO answers (
// 					id, participant_id, question_uuid, option_uuid, 
// 					is_correct, response_time_ms, points_awarded, answered_at
// 				) VALUES (
// 					$1, $2, $3, $4, $5, $6, $7, $8
// 				)
// 			`,
// 				uuid.NewString(),
// 				participantID,
// 				questionID,
// 				answer.OptionID,
// 				answer.IsCorrect,
// 				answer.ResponseTime.Milliseconds(),
// 				answer.Points,
// 				answer.AnsweredAt,
// 			)

// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return tx.Commit(ctx)
// }