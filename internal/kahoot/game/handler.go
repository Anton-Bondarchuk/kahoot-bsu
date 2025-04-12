package game

// import (
// 	"log"
// 	"net/http"
// 	"sync"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/gorilla/websocket"
// )

// // WebSocket upgrader
// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true // Allow all connections for simplicity
// 	},
// }

// // GameHandler manages game sessions and WebSocket connections
// type GameHandler struct {
// 	quizService       QuizService
// 	sessionRepository SessionRepository
// 	gameManager       *GameManager
// 	// Map of session code to active game
// 	activeSessions    map[string]*GameSession
// 	mu                sync.RWMutex
// }

// // NewGameHandler creates a new game handler instance
// func NewGameHandler(quizService QuizService, sessionRepo SessionRepository) *GameHandler {
// 	return &GameHandler{
// 		quizService:       quizService,
// 		sessionRepository: sessionRepo,
// 		gameManager:       NewGameManager(),
// 		activeSessions:    make(map[string]*GameSession),
// 	}
// }

// // RegisterRoutes registers all game-related routes
// func (h *GameHandler) RegisterRoutes(router *gin.Engine) {
// 	router.POST("/api/games", h.CreateGame)
// 	router.GET("/api/games/:code", h.GetGameStatus)
// 	router.GET("/api/ws/host/:code", h.HostGameWebSocket)
// 	router.GET("/api/ws/play/:code", h.JoinGameWebSocket)
// }

// // CreateGame handles the creation of a new game session
// func (h *GameHandler) CreateGame(c *gin.Context) {
// 	var req struct {
// 		QuizID string `json:"quiz_id" binding:"required"`
// 		HostID int64  `json:"host_id" binding:"required"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Generate a unique 6-character game code
// 	code := generateGameCode()

// 	// Retrieve quiz questions
// 	quiz, err := h.quizService.GetQuizWithQuestions(c.Request.Context(), req.QuizID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve quiz"})
// 		return
// 	}

// 	// Create new game session
// 	session := &GameSession{
// 		ID:                  uuid.NewString(),
// 		Code:                code,
// 		QuizID:              req.QuizID,
// 		HostID:              req.HostID,
// 		Quiz:                quiz,
// 		Status:              "waiting",
// 		CurrentQuestionIndex: -1, // No question active yet
// 		Players:             make(map[string]*Player),
// 		StartedAt:           time.Time{},
// 		CreatedAt:           time.Now(),
// 	}

// 	// Save to database
// 	if err := h.sessionRepository.Create(c.Request.Context(), session); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create game session"})
// 		return
// 	}

// 	// Add to active sessions
// 	h.mu.Lock()
// 	h.activeSessions[code] = session
// 	h.mu.Unlock()

// 	c.JSON(http.StatusCreated, gin.H{
// 		"game_id": session.ID,
// 		"code":    code,
// 	})
// }

// // GetGameStatus returns the current status of a game
// func (h *GameHandler) GetGameStatus(c *gin.Context) {
// 	code := c.Param("code")

// 	h.mu.RLock()
// 	session, exists := h.activeSessions[code]
// 	h.mu.RUnlock()

// 	if !exists {
// 		// Try to find from database if not in memory
// 		var err error
// 		session, err = h.sessionRepository.FindByCode(c.Request.Context(), code)
// 		if err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
// 			return
// 		}

// 		// If found in db but not in active sessions, it means it's not currently active
// 		if session.Status == "finished" {
// 			c.JSON(http.StatusOK, gin.H{
// 				"status": "finished",
// 				"game":   session,
// 			})
// 			return
// 		}

// 		c.JSON(http.StatusNotFound, gin.H{"error": "Game not active"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":      session.Status,
// 		"player_count": len(session.Players),
// 	})
// }

// // HostGameWebSocket establishes a WebSocket connection for the host
// func (h *GameHandler) HostGameWebSocket(c *gin.Context) {
// 	code := c.Param("code")

// 	h.mu.RLock()
// 	session, exists := h.activeSessions[code]
// 	h.mu.RUnlock()

// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
// 		return
// 	}

// 	// Check if the user is the host (simplified for now)
// 	// In a real app, verify the user ID from JWT token or session

// 	// Upgrade to WebSocket
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println("Error upgrading connection:", err)
// 		return
// 	}

// 	// Register host connection with game manager
// 	h.gameManager.RegisterHostConnection(code, session, conn)
// }

// // JoinGameWebSocket establishes a WebSocket connection for a player
// func (h *GameHandler) JoinGameWebSocket(c *gin.Context) {
// 	code := c.Param("code")

// 	h.mu.RLock()
// 	session, exists := h.activeSessions[code]
// 	h.mu.RUnlock()

// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
// 		return
// 	}

// 	if session.Status != "waiting" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Game already in progress"})
// 		return
// 	}

// 	// Upgrade to WebSocket
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println("Error upgrading connection:", err)
// 		return
// 	}

// 	// Wait for player to send their nickname
// 	var msg struct {
// 		Type     string `json:"type"`
// 		Nickname string `json:"nickname"`
// 		UserID   int64  `json:"user_id,omitempty"` // Optional user ID for registered users
// 	}

// 	if err := conn.ReadJSON(&msg); err != nil {
// 		log.Println("Error reading nickname:", err)
// 		conn.Close()
// 		return
// 	}

// 	if msg.Type != "join" || msg.Nickname == "" {
// 		conn.WriteJSON(map[string]string{"type": "error", "message": "Invalid join request"})
// 		conn.Close()
// 		return
// 	}

// 	// Create new player
// 	playerID := uuid.NewString()
// 	player := &Player{
// 		ID:       playerID,
// 		Nickname: msg.Nickname,
// 		UserID:   msg.UserID,
// 		Score:    0,
// 		Conn:     conn,
// 	}

// 	// Register player with game session
// 	session.mu.Lock()
// 	session.Players[playerID] = player
// 	session.mu.Unlock()

// 	// Register player connection with game manager
// 	h.gameManager.RegisterPlayerConnection(code, session, player, conn)

// 	// Send confirmation to player
// 	conn.WriteJSON(map[string]interface{}{
// 		"type":      "joined",
// 		"player_id": playerID,
// 		"game_code": code,
// 	})

// 	// Notify host about new player
// 	session.broadcastToHost(map[string]interface{}{
// 		"type":      "player_joined",
// 		"player_id": playerID,
// 		"nickname":  msg.Nickname,
// 		"players":   session.getPlayersList(),
// 	})
// }

// // generateGameCode generates a random 6-character game code
// func generateGameCode() string {
// 	// In a real implementation, ensure uniqueness and avoid confusion
// 	// (e.g., avoid similar looking characters like O and 0)
// 	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
// 	code := make([]byte, 6)
// 	for i := range code {
// 		code[i] = charset[time.Now().UnixNano()%int64(len(charset))]
// 		time.Sleep(time.Nanosecond)
// 	}
// 	return string(code)
// }