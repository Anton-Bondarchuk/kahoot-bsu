package main

import (
	"context"
	"flag"
	"kahoot_bsu/internal/infra"
	"kahoot_bsu/internal/kahoot/admin/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Version information
var (
	buildTime = "2025-04-10 19:31:25"
	buildUser = "Anton-Bondarchuk"
	version   = "1.0.0"
)

func main() {
	// Command line flags
	var (
		addr     = flag.String("addr", ":8080", "HTTP server address")
		dbURL    = flag.String("db", os.Getenv("DATABASE_URL"), "Database connection URL")
		// logLevel = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
		env      = flag.String("env", "development", "Environment (development, production)")
	)
	flag.Parse()

	// Set up logging
	log.Printf("Starting Kahoot BSU API server (version: %s, build: %s, user: %s)", 
		version, buildTime, buildUser)
	log.Printf("Environment: %s", *env)

	// Set Gin mode based on environment
	if *env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	if *dbURL == "" {
		*dbURL = "postgres://postgres:postgres@localhost:5432/postgres"
	}
	log.Printf("Connecting to database")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	db, err := pgxpool.New(ctx, *dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verify database connection
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Printf("Connected to database successfully")

	// Initialize repositories
	quizRepo := infra.NewPgQuizRepository(db)
	questionRepo := infra.NewPgQuestionRepository(db)

	// Initialize handlers
	handlers := handlers.NewHandlers(quizRepo, questionRepo)

	// Set up router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add middleware for request logging
	router.Use(gin.Logger())
	router.Use(gin.Recovery())


		// Setup template rendering
		router.LoadHTMLGlob("templates/*")
	
		// Static file serving
		router.Static("/static", "./static")
	
		// Render index.html template
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title":      "Kahoot BSU Quiz Application",
				"version":    version,
				"buildTime":  buildTime,
				"buildUser":  buildUser,
				"serverTime": time.Now().Format("2006-01-02 15:04:05"),
			})
		})
	

	api := router.Group("/api")
	{
		// Quiz routes
		api.GET("/quizzes", handlers.GetUserQuizzes)
		api.POST("/quizzes", handlers.CreateQuiz)
		api.GET("/quizzes/:id", handlers.GetQuiz)
		api.PUT("/quizzes/:id", handlers.UpdateQuiz)
		api.DELETE("/quizzes/:id", handlers.DeleteQuiz)

		// Question routes
		api.GET("/quizzes/:id/questions", handlers.GetQuizQuestions)
		api.POST("/quizzes/:id/questions", handlers.AddQuizQuestion)
		api.PUT("/questions/:question_id", handlers.UpdateQuestion)
		api.DELETE("/questions/:question_id", handlers.DeleteQuestion)
	}

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "up",
			"version":   version,
			"buildTime": buildTime,
			"buildUser": buildUser,
		})
	})

	// Set up HTTP server
	srv := &http.Server{
		Addr:    *addr,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on %s", *addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
