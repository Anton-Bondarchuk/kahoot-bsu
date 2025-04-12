package infra

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log/slog"
// 	"net/url"
// 	"os"
// 	"time"

// 	"github.com/ThreeDotsLabs/watermill"
// 	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
// 	"github.com/ThreeDotsLabs/watermill/message"
// 	"github.com/google/uuid"
// )


// const  (
// 	GameEventsTopic    = "game.events"
// 	PlayerEventsTopic  = "player.events"
// 	AnswerEventsTopic  = "answer.events"
// 	QuestionEventsTopic = "question.events"
// )

// type amqpMessagesPublisher struct {
// 	pub *amqp.Publisher
// }


// func NewAmqpMessagesPublisher() (*amqpMessagesPublisher, <-chan *message.Message, func() error) {
// 	uri := os.Getenv("")
// 	if uri == "" {
// 		uri = "amqp://guest:guest@localhost:5672/"
// 	}
// 	amqpConfig := amqp.NewDurableQueueConfig(uri)


// 	wmLogger := watermill.NewStdLogger(false, false)
// 	publisher, err := amqp.NewPublisher(amqpConfig, wmLogger)
// 	if err != nil {
// 		panic(err)
// 	}
// 	sub, err = amqp.NewSubscriber(amqpConfig, wmLogger)
// 	if (err != nil) {
// 		panic(err)
// 	}



// }
