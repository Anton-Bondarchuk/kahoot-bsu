package telegram

import (
	"context"
	"kahoot_bsu/internal/config"
	"kahoot_bsu/internal/domain/models"
	"kahoot_bsu/internal/infra/clients"
	infra "kahoot_bsu/internal/infra/persistence"
	"kahoot_bsu/internal/infra/services"
	"time"

	"kahoot_bsu/internal/logger/handlers/slogpretty"
	"kahoot_bsu/internal/service/fsm"
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
	Config *config.Config
	Conn   *pgxpool.Pool
	Bot    *models.Bot
	Log    *slog.Logger
	router *fsm.Router
}

func NewAppTelegram() (
	app *AppTelegram,
	close func() error,
) {

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	telegramBot := newBot(cfg.BotConfig)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := newPgxConn(ctx, cfg.StorageConfig)

	log.Info("Connected to database successfully")

	ctx = context.Background()

	redisStorage := NewRedisStorage(ctx, cfg.RedisConfig)

	router := fsm.NewRouter(redisStorage)

	emailClient := clients.NewEmailClient(cfg.EmailConfig)
	emailService := services.NewEmailService(emailClient)
	userRepo := infra.NewPgUserRepository(db)
	otpGenerator := services.NewVerificationOTPGenerator(6)

	fSMHandler := messages.NewFSMHandler(emailService, userRepo, otpGenerator)

	router.Register(fsm.StateAwaitingLogin, fSMHandler.HandleLogin)
	router.Register(fsm.StateAwaitingOTP, fSMHandler.HandleOTP)
	router.Register(fsm.StateRegistered, fSMHandler.HandleRegistered)

	app = &AppTelegram{
		Config: cfg,
		Bot:    telegramBot,
		Conn:   db,
		Log:    log,
		router: router,
	}

	closeFunc := func() error {
		var err error

		// if i will have errors I should use it
		return err
	}

	return app, closeFunc
}

func newBot(cfg config.BotConfig) *models.Bot {
	botAPI, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		panic(err)
	}

	botAPI.Debug = cfg.Debug

	// Note: add With and functional option pattern
	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Timeout

	updates := botAPI.GetUpdatesChan(u)

	return &models.Bot{
		Telegram:      botAPI,
		UpdateChannel: updates,
	}
}

func newPgxConn(ctx context.Context, cfg config.StorageConfig) *pgxpool.Pool {
	db, err := pgxpool.New(ctx, cfg.DatabaseUrl)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		panic(err)
	}

	return db
}

func NewRedisStorage(ctx context.Context, cfg config.RedisConfig) *infra.RedisStorage {
	storage := infra.NewRedisStorage(cfg)

	if err := storage.Ping(context.Background()); err != nil {
		panic(err)
	}

	return storage
}

func Start(ctx context.Context, a *AppTelegram) {
	for update := range a.Bot.UpdateChannel {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID
		message := update.Message
		states := 
	
		fsm := fsm.NewFSMContext(ctx, a.router.Storage, chatID, userID)

		if update.Message.IsCommand() {
			// handleCommand(a, message, fsm)
		} else {
			if err := a.router.ProcessUpdate(ctx, message, a.Bot, fsm); err != nil {
				a.Log.Error("Error processing update: %v", err)
			}

		}
	}
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
