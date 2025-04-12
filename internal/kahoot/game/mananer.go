package game

// import (
// 	"encoding/json"
// 	"log"
// 	"sync"

// 	"github.com/gorilla/websocket"
// )

// // GameManager handles the real-time game logic and WebSocket communication
// type GameManager struct {
// 	// Map of session code to host connection
// 	hostConns map[string]*websocket.Conn
// 	// Map of session code to map of player ID to player connection
// 	playerConns map[string]map[string]*websocket.Conn
// 	mu          sync.RWMutex
// }

// // NewGameManager creates a new game manager
// func NewGameManager() *GameManager {
// 	return &GameManager{
// 		hostConns:   make(map[string]*websocket.Conn),
// 		playerConns: make(map[string]map[string]*websocket.Conn),
// 	}
// }

// // RegisterHostConnection registers a host WebSocket connection
// func (gm *GameManager) RegisterHostConnection(code string, session *GameSession, conn *websocket.Conn) {
// 	gm.mu.Lock()
// 	gm.hostConns[code] = conn
// 	gm.mu.Unlock()

// 	// Start goroutine to handle host messages
// 	go gm.handleHostMessages(code, session, conn)
// }

// // RegisterPlayerConnection registers a player WebSocket connection
// func (gm *GameManager) RegisterPlayerConnection(code string, session *GameSession, player *Player, conn *websocket.Conn) {
// 	gm.mu.Lock()
// 	if _, exists := gm.playerConns[code]; !exists {
// 		gm.playerConns[code] = make(map[string]*websocket.Conn)
// 	}
// 	gm.playerConns[code][player.ID] = conn
// 	gm.mu.Unlock()

// 	// Start goroutine to handle player messages
// 	go gm.handlePlayerMessages(code, session, player, conn)
// }

// // handleHostMessages processes messages from the host
// func (gm *GameManager) handleHostMessages(code string, session *GameSession, conn *websocket.Conn) {
// 	defer func() {
// 		conn.Close()
// 		gm.mu.Lock()
// 		delete(gm.hostConns, code)
// 		gm.mu.Unlock()
// 		log.Printf("Host disconnected from game %s\n", code)
// 	}()

// 	for {
// 		var msg map[string]interface{}
// 		if err := conn.ReadJSON(&msg); err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("Host connection error: %v\n", err)
// 			}
// 			break
// 		}

// 		// Process host messages
// 		msgType, ok := msg["type"].(string)
// 		if !ok {
// 			continue
// 		}

// 		switch msgType {
// 		case "start_game":
// 			// Start the game
// 			log.Printf("Starting game %s\n", code)
// 			session.startGame()

// 		case "next_question":
// 			// Move to next question
// 			log.Printf("Moving to next question in game %s\n", code)
// 			session.nextQuestion()

// 		case "end_game":
// 			// End the game
// 			log.Printf("Ending game %s\n", code)
// 			session.endGame()
// 		}
// 	}
// }

// // handlePlayerMessages processes messages from players
// func (gm *GameManager) handlePlayerMessages(code string, session *GameSession, player *Player, conn *websocket.Conn) {
// 	defer func() {
// 		conn.Close()
		
// 		// Remove player from game session
// 		session.mu.Lock()
// 		delete(session.Players, player.ID)
// 		session.mu.Unlock()
		
// 		// Remove connection from manager
// 		gm.mu.Lock()
// 		if conns, exists := gm.playerConns[code]; exists {
// 			delete(conns, player.ID)
// 		}
// 		gm.mu.Unlock()
		
// 		// Notify host about player leaving
// 		session.broadcastToHost(map[string]interface{}{
// 			"type":      "player_left",
// 			"player_id": player.ID,
// 			"players":   session.getPlayersList(),
// 		})
		
// 		log.Printf("Player %s disconnected from game %s\n", player.Nickname, code)
// 	}()

// 	for {
// 		var msg map[string]interface{}
// 		if err := conn.ReadJSON(&msg); err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("Player connection error: %v\n", err)
// 			}
// 			break
// 		}

// 		// Process player messages
// 		msgType, ok := msg["type"].(string)
// 		if !ok {
// 			continue
// 		}

// 		switch msgType {
// 		case "answer":
// 			// Process player answer
// 			if session.Status != "question_active" {
// 				continue
// 			}

// 			optionIDRaw, ok := msg["option_id"].(string)
// 			if !ok {
// 				continue
// 			}

// 			log.Printf("Player %s answered in game %s\n", player.Nickname, code)
// 			session.processAnswer(player.ID, optionIDRaw)
// 		}
// 	}
// }

// // broadcastToAll sends a message to all connected players and the host
// func (gm *GameManager) broadcastToAll(code string, message interface{}) {
// 	jsonMessage, err := json.Marshal(message)
// 	if err != nil {
// 		log.Printf("Error marshaling broadcast message: %v\n", err)
// 		return
// 	}

// 	// Send to host
// 	gm.mu.RLock()
// 	hostConn, hostExists := gm.hostConns[code]
// 	gm.mu.RUnlock()

// 	if hostExists {
// 		if err := hostConn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
// 			log.Printf("Error sending message to host: %v\n", err)
// 		}
// 	}

// 	// Send to all players
// 	gm.mu.RLock()
// 	playerConnsMap, exists := gm.playerConns[code]
// 	playerConns := make([]*websocket.Conn, 0, len(playerConnsMap))
// 	if exists {
// 		for _, conn := range playerConnsMap {
// 			playerConns = append(playerConns, conn)
// 		}
// 	}
// 	gm.mu.RUnlock()

// 	for _, conn := range playerConns {
// 		if err := conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
// 			log.Printf("Error sending message to player: %v\n", err)
// 			// In a production system, we'd handle connection errors more gracefully
// 		}
// 	}
// }