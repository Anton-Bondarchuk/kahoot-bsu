package infra

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log/slog"
// 	"os"
// 	"time"

// 	"github.com/ThreeDotsLabs/watermill"
// 	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
// 	"github.com/ThreeDotsLabs/watermill/message"
// 	"github.com/google/uuid"
// )

// // Topic constants for different event types
// const (
// 	GameEventsTopic    = "game.events"
// 	PlayerEventsTopic  = "player.events"
// 	AnswerEventsTopic  = "answer.events"
// 	QuestionEventsTopic = "question.events"
// )

// // Publisher handles publishing messages to RabbitMQ using Watermill
// type Publisher struct {
// 	pub     *amqp.Publisher
// 	log     slog.Logger
// }

// // NewPublisher creates a new RabbitMQ publisher using Watermill
// func NewPublisher(log slog.Logger) (*Publisher, <-chan *message.Message, func() error, error) {
// 	// Get RabbitMQ URI from environment or use default
// 	uri := os.Getenv("AMQP_URI")
// 	if uri == "" {
// 		uri = "amqp://guest:guest@localhost:5672/"
// 	}
	
// 	// Configure AMQP connection for Watermill
// 	amqpConfig := amqp.NewDurableQueueConfig(uri)
	
// 	// Create Watermill logger
// 	wmLogger := watermill.NewStdLogger(false, false)
	
// 	// Create Watermill AMQP publisher
// 	publisher, err := amqp.NewPublisher(amqpConfig, wmLogger)
// 	if err != nil {
// 		return nil, nil, nil, fmt.Errorf("failed to create AMQP publisher: %w", err)
// 	}
	
// 	// Create subscriber for receiving messages
// 	subscriber, err := amqp.NewSubscriber(amqpConfig, wmLogger)
// 	if err != nil {
// 		return nil, nil, nil, fmt.Errorf("failed to create AMQP subscriber: %w", err)
// 	}
	
// 	// Subscribe to events we're interested in
// 	messages, err := subscriber.Subscribe(context.Background(), GameEventsTopic)
// 	if err != nil {
// 		return nil, nil, nil, fmt.Errorf("failed to subscribe to topics: %w", err)
// 	}
	
// 	// Create close function
// 	closeFunc := func() error {
// 		err := publisher.Close()
// 		if err != nil {
// 			return fmt.Errorf("failed to close publisher: %w", err)
// 		}
// 		return nil
// 	}
	
// 	return &Publisher{
// 		pub: publisher,
// 		log: log,
// 	}, messages, closeFunc, nil
// }

// // PublishEvent publishes an event to the specified topic
// func (p *Publisher) PublishEvent(ctx context.Context, topic string, eventType string, payload interface{}) error {
// 	// Create message with event type and payload
// 	msg := models.Message{
// 		Type:      eventType,
// 		Timestamp: time.Now(),
// 		Payload:   payload,
// 	}
	
// 	// Marshal message to JSON
// 	data, err := json.Marshal(msg)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal message: %w", err)
// 	}
	
// 	// Create Watermill message
// 	watermillMsg := message.NewMessage(watermill.NewUUID(), data)
	
// 	// Add metadata
// 	watermillMsg.Metadata.Set("event_type", eventType)
	
// 	// Publish message
// 	if err := p.pub.Publish(topic, watermillMsg); err != nil {
// 		return fmt.Errorf("failed to publish message: %w", err)
// 	}
	
// 	p.log.Info("Published %s event to %s topic", eventType, topic)
// 	return nil
// }

// // PublishGameEvent publishes a game-related event
// func (p *Publisher) PublishGameEvent(ctx context.Context, eventType string, payload interface{}) error {
// 	return p.PublishEvent(ctx, GameEventsTopic, eventType, payload)
// }

// // PublishPlayerEvent publishes a player-related event
// func (p *Publisher) PublishPlayerEvent(ctx context.Context, eventType string, payload interface{}) error {
// 	return p.PublishEvent(ctx, PlayerEventsTopic, eventType, payload)
// }

// // PublishAnswerEvent publishes an answer-related event
// func (p *Publisher) PublishAnswerEvent(ctx context.Context, eventType string, payload interface{}) error {
// 	return p.PublishEvent(ctx, AnswerEventsTopic, eventType, payload)
// }

// // PublishQuestionEvent publishes a question-related event
// func (p *Publisher) PublishQuestionEvent(ctx context.Context, eventType string, payload interface{}) error {
// 	return p.PublishEvent(ctx, QuestionEventsTopic, eventType, payload)
// }

// // Event-specific methods for easier use

// // PublishPlayerJoined publishes a player joined event
// func (p *Publisher) PublishPlayerJoined(ctx context.Context, player *models.Player, gameID uuid.UUID) error {
// 	payload := map[string]interface{}{
// 		"player_id": player.ID,
// 		"game_id":   gameID,
// 		"username":  player.Username,
// 	}
// 	return p.PublishPlayerEvent(ctx, models.EventPlayerJoined, payload)
// }

// // PublishPlayerLeft publishes a player left event
// func (p *Publisher) PublishPlayerLeft(ctx context.Context, playerID, gameID uuid.UUID) error {
// 	payload := map[string]interface{}{
// 		"player_id": playerID,
// 		"game_id":   gameID,
// 	}
// 	return p.PublishPlayerEvent(ctx, models.EventPlayerLeft, payload)
// }

// // PublishGameCreated publishes a game created event
// func (p *Publisher) PublishGameCreated(ctx context.Context, game *models.Game) error {
// 	return p.PublishGameEvent(ctx, models.EventGameCreated, game)
// }

// // PublishGameStarted publishes a game started event
// func (p *Publisher) PublishGameStarted(ctx context.Context, game *models.Game) error {
// 	return p.PublishGameEvent(ctx, models.EventGameStarted, game)
// }

// // PublishGameEnded publishes a game ended event
// func (p *Publisher) PublishGameEnded(ctx context.Context, game *models.Game) error {
// 	return p.PublishGameEvent(ctx, models.EventGameEnded, game)
// }

// // PublishQuestionStarted publishes a question started event
// func (p *Publisher) PublishQuestionStarted(ctx context.Context, game *models.Game, questionIndex int) error {
// 	question := game.Questions[questionIndex]
// 	payload := map[string]interface{}{
// 		"game_id":     game.ID,
// 		"question_id": question.ID,
// 		"index":       questionIndex,
// 		"time_limit":  question.TimeLimit,
// 	}
// 	return p.PublishQuestionEvent(ctx, models.EventQuestionStarted, payload)
// }

// // PublishQuestionEnded publishes a question ended event
// func (p *Publisher) PublishQuestionEnded(ctx context.Context, game *models.Game, questionIndex int) error {
// 	question := game.Questions[questionIndex]
// 	payload := map[string]interface{}{
// 		"game_id":     game.ID,
// 		"question_id": question.ID,
// 		"index":       questionIndex,
// 	}
// 	return p.PublishQuestionEvent(ctx, models.EventQuestionEnded, payload)
// }

// // PublishAnswerSubmitted publishes an answer submitted event
// func (p *Publisher) PublishAnswerSubmitted(ctx context.Context, answer *models.Answer) error {
// 	return p.PublishAnswerEvent(ctx, models.EventAnswerSubmitted, answer)
// }

// // PublishLeaderboardUpdated publishes a leaderboard updated event
// func (p *Publisher) PublishLeaderboardUpdated(ctx context.Context, gameID uuid.UUID, leaderboard interface{}) error {
// 	payload := map[string]interface{}{
// 		"game_id":     gameID,
// 		"leaderboard": leaderboard,
// 	}
// 	return p.PublishGameEvent(ctx, models.EventLeaderboardUpdated, payload)
// }