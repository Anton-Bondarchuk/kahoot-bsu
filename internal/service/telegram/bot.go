package telegram

import (
	"context"
	"errors"
	"fmt"
	"kahoot_bsu/internal/app/command"
	"kahoot_bsu/internal/app/messages"
	"kahoot_bsu/internal/config"
	"kahoot_bsu/internal/domain/models"
	"kahoot_bsu/internal/logger/handlers/slogpretty"
	"log/slog"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type AppTelegram struct {
	Conn *pgxpool.Pool
	Bot  *models.Bot
	Log  *slog.Logger
}

func NewAppTelegram() (
	app *AppTelegram,
	close func() error,
) {
	
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	telegramBot, err := newBot(cfg.BotConfig)
	if err != nil {
		return nil, func() error {
			return fmt.Errorf("failed to initialize bot: %w", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, cfg.StorageConfig.DatabaseUrl)
	if err != nil {
		return nil, func() error {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
	}

	if err := db.Ping(ctx); err != nil {
		return nil, func() error {
			db.Close()
			return errors.Join(
				fmt.Errorf("failed to ping database: %w", err),
			)
		}
	}
	
	log.Info("Connected to database successfully")

	app = &AppTelegram{
		Bot:    telegramBot,
		Conn:     db,
		Log: log,
	}

	closeFunc := func() error {
		var err error
		
		// if i will have errors I should use it
		return err
	}

	return app, closeFunc
}


func newBot(cfg config.BotConfig) (*models.Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	botAPI.Debug = cfg.Debug

	// TODO: add With and functional option pattern
	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Timeout

	updates := botAPI.GetUpdatesChan(u)

	return &models.Bot{
		Telegram:      botAPI,
		UpdateChannel: updates}, nil
}

func Start(b *models.Bot) {
	for update := range b.UpdateChannel {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			handleCommand(b, update.Message)
		} else {
			handleMessages(b, update.Message)
		}
	}
}
type CommandInterface interface {
	Execute(message *tgbotapi.Message)
}

func handleCommand(b *models.Bot, message *tgbotapi.Message) {
	comandHandler := command.New(b, "https://af09-185-53-133-77.ngrok-free.app/")
	// vericationRepo := infra.NewPgVerificationCodeRepository()

	commandStrategy := map[string]CommandInterface{
		"start":    &command.StartCommand{CommandHandler: comandHandler},
		"register": &command.RegisterCommand{CommandHandler: comandHandler},
		"kahoot":   &command.KahootComand{CommandHandler: comandHandler},
		"help":     &command.HelpCommand{CommandHandler: comandHandler},
		"unknown":  &command.UnknownCommand{CommandHandler: comandHandler},
	}

	_, ok := commandStrategy[message.Command()]

	if !ok {
		commandStrategy["unknown"].Execute(message)
		return
	}

	commandStrategy[message.Command()].Execute(message)
}

func handleMessages(b *models.Bot, message *tgbotapi.Message) {
	messHandler := messages.New(b)

	emailReg := &messages.HandleEmailRegistrationMessenger{MessagesHandler: messHandler}
	emailReg.Execute(message)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}


func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}