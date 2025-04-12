package main

import (
	// "context"
	// "kahoot_bsu/internal/config"
	// "kahoot_bsu/internal/logger/handlers/slogpretty"
	// "log/slog"
	// "os"
	// "os/signal"
	// "syscall"
	// "time"

	// "kahoot/players-service/internal/rabbitmq"
	// "kahoot/players-service/internal/repository"
	// "kahoot/players-service/internal/server"
	// "kahoot/players-service/internal/service"
	// "kahoot/players-service/pkg/logger"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)


// func main() {
// 	// Initialize logger
// 	log := setupLogger(cnf.Env)

// 	// Load configuration
// 	cfg := config.MustLoad()

// 	// Connect to database
// 	repo, err := repository.NewRepository(cfg.Database)
// 	if err != nil {
// 		log.Error("failed to initialize repository: %v", err)
// 	}
// 	defer repo.Close()

// 	// Connect to RabbitMQ
// 	publisher, err := rabbitmq.NewPublisher(cfg.RabbitMQ)
// 	if err != nil {
// 		log.Error("Failed to connect to RabbitMQ: %v", err)
// 	}
// 	defer publisher.Close()

// 	// Initialize services
// 	gameService := service.NewGameService(repo, publisher, log)
// 	playerService := service.NewPlayerService(repo, publisher, log)

// 	// Create and start HTTP server
// 	srv := server.NewServer(cfg, gameService, playerService, log)
// 	go func() {
// 		if err := srv.Start(); err != nil {
// 			log.Error("Failed to start server: %v", err)
// 		}
// 	}()

// 	// Graceful shutdown
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit

// 	log.Info("Shutting down server...")
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	if err := srv.Shutdown(ctx); err != nil {
// 		log.Error("Server forced to shutdown: %v", err)
// 	}

// 	log.Info("Server exited properly")
// }

// func setupLogger(env string) *slog.Logger {
// 	var log *slog.Logger

// 	switch env {
// 	case envLocal:
// 		log = setupPrettySlog()
// 	case envDev:
// 		log = slog.New(
// 			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
// 		)
// 	case envProd:
// 		log = slog.New(
// 			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
// 		)
// 	}

// 	return log
// }


// func setupPrettySlog() *slog.Logger {
// 	opts := slogpretty.PrettyHandlerOptions{
// 		SlogOpts: &slog.HandlerOptions{
// 			Level: slog.LevelDebug,
// 		},
// 	}

// 	handler := opts.NewPrettyHandler(os.Stdout)

// 	return slog.New(handler)
// }