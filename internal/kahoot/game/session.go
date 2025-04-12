package game

// import (
// 	"sort"
// 	"sync"
// 	"time"

// 	"github.com/gorilla/websocket"
// )

// // GameSession represents an active quiz game session
// type GameSession struct {
// 	ID                  string
// 	Code                string
// 	QuizID              string
// 	HostID              int64
// 	Quiz                *Quiz
// 	Status              string // "waiting", "question_active", "reviewing", "finished"
// 	CurrentQuestionIndex int
// 	Players             map[string]*Player
// 	StartedAt           time.Time
// 	EndedAt             time.Time
// 	CreatedAt           time.Time
// 	mu                  sync.RWMutex
// 	hostConn            *websocket.Conn // Connection to the host
// }

// // Player represents a participant in a game session
// type Player struct {
// 	ID          string
// 	Nickname    string
// 	UserID      int64 // Optional, for registered users
// 	Score       int
// 	Conn        *websocket.Conn
// 	Answers     map[string]PlayerAnswer // Map of question ID to answer
// }

// // PlayerAnswer represents a player's answer to a question
// type PlayerAnswer struct {
// 	OptionID     string
// 	AnsweredAt   time.Time
// 	ResponseTime time.Duration
// 	IsCorrect    bool
// 	Points       int
// }

// // startGame initiates the game
// func (gs *GameSession) startGame() {
// 	gs.mu.Lock()
// 	if gs.Status != "waiting" {
// 		gs.mu.Unlock()
// 		return
// 	}
// 	gs.Status = "started"
// 	gs.StartedAt = time.Now()
// 	gs.mu.Unlock()

// 	// Broadcast game started to all players
// 	gs.broadcastToAll(map[string]interface{}{
// 		"type":      "game_started",
// 		"game_code": gs.Code,
// 		"quiz_name": gs.Quiz.Title,
// 	})

// 	// Start the first question after a short delay
// 	time.Sleep(3 * time.Second)
// 	gs.nextQuestion()
// }

// // nextQuestion moves to the next question or ends the game if all questions are answered
// func (gs *GameSession) nextQuestion() {
// 	gs.mu.Lock()
	
// 	// Check if the game is in progress
// 	if gs.Status == "finished" {
// 		gs.mu.Unlock()
// 		return
// 	}
	
// 	// Increment question index
// 	gs.CurrentQuestionIndex++
	
// 	// Check if we've reached the end of the quiz
// 	if gs.CurrentQuestionIndex >= len(gs.Quiz.Questions) {
// 		gs.Status = "finished"
// 		gs.EndedAt = time.Now()
// 		gs.mu.Unlock()
		
// 		// End the game
// 		gs.endGame()
// 		return
// 	}
	
// 	// Get the current question
// 	currentQuestion := gs.Quiz.Questions[gs.CurrentQuestionIndex]
// 	gs.Status = "question_active"
// 	questionStartTime := time.Now()
// 	gs.mu.Unlock()
	
// 	// Prepare question data for the clients
// 	questionData := map[string]interface{}{
// 		"type":         "question",
// 		"question_idx": gs.CurrentQuestionIndex + 1, // 1-based for display
// 		"question_id":  currentQuestion.UUID,
// 		"text":         currentQuestion.Text,
// 		"time_limit":   currentQuestion.TimeLimit,
// 		"options":      prepareOptionsForClient(currentQuestion.Options),
// 		"total_questions": len(gs.Quiz.Questions),
// 	}
	
// 	// Send question to all players
// 	gs.broadcastToAll(questionData)
	
// 	// Wait for the question's time limit
// 	time.Sleep(time.Duration(currentQuestion.TimeLimit) * time.Second)
	
// 	// End the question period
// 	gs.mu.Lock()
// 	if gs.Status == "question_active" && gs.CurrentQuestionIndex < len(gs.Quiz.Questions) {
// 		gs.Status = "reviewing"
// 	}
// 	gs.mu.Unlock()
	
// 	// Calculate results for this question
// 	results := gs.calculateQuestionResults(currentQuestion.UUID)
	
// 	// Send results to all players
// 	gs.broadcastToAll(map[string]interface{}{
// 		"type":     "question_results",
// 		"question_id": currentQuestion.UUID,
// 		"correct_option": getCorrectOptionID(currentQuestion.Options),
// 		"results":  results,
// 	})
	
// 	// Wait for a short period to let players see the results
// 	time.Sleep(5 * time.Second)
// }

// // prepareOptionsForClient removes the "is_correct" flag from options before sending to clients
// func prepareOptionsForClient(options []Option) []map[string]interface{} {
// 	result := make([]map[string]interface{}, len(options))
// 	for i, opt := range options {
// 		result[i] = map[string]interface{}{
// 			"id":   opt.UUID,
// 			"text": opt.Text,
// 		}
// 	}
// 	return result
// }

// // getCorrectOptionID returns the ID of the correct option
// func getCorrectOptionID(options []Option) string {
// 	for _, opt := range options {
// 		if opt.IsCorrect {
// 			return opt.UUID
// 		}
// 	}
// 	return ""
// }

// // processAnswer handles a player's answer submission
// func (gs *GameSession) processAnswer(playerID, optionID string) {
// 	gs.mu.Lock()
// 	defer gs.mu.Unlock()
	
// 	// Verify the game state
// 	if gs.Status != "question_active" || gs.CurrentQuestionIndex >= len(gs.Quiz.Questions) {
// 		return
// 	}
	
// 	// Get the player
// 	player, exists := gs.Players[playerID]
// 	if !exists {
// 		return
// 	}
	
// 	// Get the current question
// 	currentQuestion := gs.Quiz.Questions[gs.CurrentQuestionIndex]
	
// 	// Check if player has already answered this question
// 	if player.Answers == nil {
// 		player.Answers = make(map[string]PlayerAnswer)
// 	}
// 	if _, answered := player.Answers[currentQuestion.UUID]; answered {
// 		return // Player already answered
// 	}
	
// 	// Record the answer time
// 	now := time.Now()
// 	responseTime := now.Sub(gs.StartedAt)
	
// 	// Find the selected option and check if it's correct
// 	var isCorrect bool
// 	for _, opt := range currentQuestion.Options {
// 		if opt.UUID == optionID {
// 			isCorrect = opt.IsCorrect
// 			break
// 		}
// 	}
	
// 	// Calculate points based on correctness and response time
// 	points := 0
// 	if isCorrect {
// 		// Formula: Base points * (1 - (responseTime / timeLimit))
// 		// This rewards faster responses
// 		timeLimit := time.Duration(currentQuestion.TimeLimit) * time.Second
// 		timeRatio := float64(responseTime) / float64(timeLimit)
// 		if timeRatio > 1 {
// 			timeRatio = 1 // Cap at 1 for late responses
// 		}
// 		points = int(float64(currentQuestion.Points) * (1 - timeRatio * 0.5)) // Min 50% for correct answers
// 	}
	
// 	// Record the answer
// 	player.Answers[currentQuestion.UUID] = PlayerAnswer{
// 		OptionID:     optionID,
// 		AnsweredAt:   now,
// 		ResponseTime: responseTime,
// 		IsCorrect:    isCorrect,
// 		Points:       points,
// 	}
	
// 	// Update the player's score
// 	player.Score += points
	
// 	// Notify the host about the answer
// 	gs.broadcastToHost(map[string]interface{}{
// 		"type":        "player_answered",
// 		"player_id":   playerID,
// 		"nickname":    player.Nickname,
// 		"question_id": currentQuestion.UUID,
// 		"is_correct":  isCorrect,
// 		"points":      points,
// 	})
// }

// // calculateQuestionResults compiles the results for a question
// func (gs *GameSession) calculateQuestionResults(questionID string) map[string]interface{} {
// 	gs.mu.RLock()
// 	defer gs.mu.RUnlock()
	
// 	// Collect answers statistics
// 	totalPlayers := len(gs.Players)
// 	answered := 0
// 	correct := 0
// 	optionCounts := make(map[string]int)
	
// 	// Process each player's answer
// 	for _, player := range gs.Players {
// 		if answer, exists := player.Answers[questionID]; exists {
// 			answered++
// 			if answer.IsCorrect {
// 				correct++
// 			}
// 			optionCounts[answer.OptionID]++
// 		}
// 	}
	
// 	// Get current leaderboard
// 	leaderboard := gs.getLeaderboard()
	
// 	return map[string]interface{}{
// 		"total_players": totalPlayers,
// 		"answered":      answered,
// 		"correct":       correct,
// 		"option_counts": optionCounts,
// 		"leaderboard":   leaderboard,
// 	}
// }

// // endGame finalizes the game and sends final results
// func (gs *GameSession) endGame() {
// 	gs.mu.Lock()
// 	gs.Status = "finished"
// 	gs.EndedAt = time.Now()
// 	leaderboard := gs.getLeaderboard()
// 	gs.mu.Unlock()
	
// 	// Send final results to all players
// 	gs.broadcastToAll(map[string]interface{}{
// 		"type":        "game_ended",
// 		"leaderboard": leaderboard,
// 		"game_code":   gs.Code,
// 		"quiz_name":   gs.Quiz.Title,
// 	})
// }

// // getPlayersList returns a list of players for the client
// func (gs *GameSession) getPlayersList() []map[string]interface{} {
// 	gs.mu.RLock()
// 	defer gs.mu.RUnlock()
	
// 	players := make([]map[string]interface{}, 0, len(gs.Players))
// 	for id, player := range gs.Players {
// 		players = append(players, map[string]interface{}{
// 			"id":       id,
// 			"nickname": player.Nickname,
// 			"score":    player.Score,
// 		})
// 	}
	
// 	return players
// }

// // getLeaderboard returns the current leaderboard
// func (gs *GameSession) getLeaderboard() []map[string]interface{} {
// 	// Create a slice of players for sorting
// 	players := make([]*Player, 0, len(gs.Players))
// 	for _, player := range gs.Players {
// 		players = append(players, player)
// 	}
	
// 	// Sort by score (descending)
// 	sort.Slice(players, func(i, j int) bool {
// 		return players[i].Score > players[j].Score
// 	})
	
// 	// Create the leaderboard
// 	leaderboard := make([]map[string]interface{}, len(players))
// 	for i, player := range players {
// 		leaderboard[i] = map[string]interface{}{
// 			"position": i + 1,
// 			"id":       player.ID,
// 			"nickname": player.Nickname,
// 			"score":    player.Score,
// 		}
// 	}
	
// 	return leaderboard
// }

// // broadcastToAll sends a message to all players and the host
// func (gs *GameSession) broadcastToAll(message interface{}) {
// 	// This is delegated to the GameManager in a real implementation
// 	// Here we'll use a simple approach for demonstration
// 	gs.mu.RLock()
// 	defer gs.mu.RUnlock()
	
// 	// Send to host
// 	if gs.hostConn != nil {
// 		gs.hostConn.WriteJSON(message)
// 	}
	
// 	// Send to all players
// 	for _, player := range gs.Players {
// 		if player.Conn != nil {
// 			player.Conn.WriteJSON(message)
// 		}
// 	}
// }

// // broadcastToHost sends a message to the host only
// func (gs *GameSession) broadcastToHost(message interface{}) {
// 	gs.mu.RLock()
// 	defer gs.mu.RUnlock()
	
// 	if gs.hostConn != nil {
// 		gs.hostConn.WriteJSON(message)
// 	}
// }