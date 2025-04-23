package fsm

import (
	"context"
	"errors"
	"slices"

	"kahoot_bsu/internal/domain/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Fix me

// Storage interface for persisting FSM states
type Storage interface {
	// Get current state for a user in a specific chat
	Get(ctx context.Context, chatID int64, userID int64) (models.State, error)

	// Set state for a user in a specific chat
	Set(ctx context.Context, chatID int64, userID int64, state models.State) error

	// Remove state for a user in a specific chat
	Delete(ctx context.Context, chatID int64, userID int64) error

	// Get data associated with the current state
	GetData(ctx context.Context, chatID int64, userID int64, key string) (interface{}, error)

	// Set data for the current state
	SetData(ctx context.Context, chatID int64, userID int64, key string, value interface{}) error

	// Clear all data for the current state
	ClearData(ctx context.Context, chatID int64, userID int64) error
}

// FSMContext manages the state for a specific user in a specific chat
type FSMContext struct {
	storage Storage
	chatID  int64
	userID  int64
	ctx     context.Context
}

// NewFSMContext creates a new FSM context for a user in a chat
func NewFSMContext(ctx context.Context, storage Storage, chatID, userID int64) *FSMContext {
	return &FSMContext{
		storage: storage,
		chatID:  chatID,
		userID:  userID,
		ctx:     ctx,
	}
}

// Current gets the current state
func (f *FSMContext) Current() (models.State, error) {
	return f.storage.Get(f.ctx, f.chatID, f.userID)
}

// IsInState checks if the current state matches any of the provided states
func (f *FSMContext) IsInState(states ...models.State) (bool, error) {
	current, err := f.Current()
	if err != nil {
		return false, err
	}

	if slices.Contains(states, current) {
		return true, nil
	}

	return false, nil
}

// Set sets a new state
func (f *FSMContext) Set(state models.State) error {
	return f.storage.Set(f.ctx, f.chatID, f.userID, state)
}

// Finish resets the state to default (ends the conversation)
func (f *FSMContext) Finish() error {
	return f.storage.Delete(f.ctx, f.chatID, f.userID)
}

// ResetState resets the state but keeps the data
func (f *FSMContext) ResetState() error {
	return f.storage.Set(f.ctx, f.chatID, f.userID, models.DefaultState)
}

// GetData gets data associated with current state
func (f *FSMContext) GetData(key string) (interface{}, error) {
	return f.storage.GetData(f.ctx, f.chatID, f.userID, key)
}

// SetData sets data for the current state
func (f *FSMContext) SetData(key string, value interface{}) error {
	return f.storage.SetData(f.ctx, f.chatID, f.userID, key, value)
}

// UpdateData updates data if it exists, otherwise sets it
func (f *FSMContext) UpdateData(key string, updateFn func(any) interface{}) error {
	value, err := f.GetData(key)
	if err != nil {
		return err
	}

	newValue := updateFn(value)
	return f.SetData(key, newValue)
}

// ClearData removes all data associated with the state
func (f *FSMContext) ClearData() error {
	return f.storage.ClearData(f.ctx, f.chatID, f.userID)
}

// HandlerFunc is a function that handles a message in a specific state
type HandlerFunc func(ctx context.Context, fsm *FSMContext, message *tgbotapi.Message, bot *models.Bot) error

// Router manages state transitions and handlers (similar to aiogram's Router)
type Router struct {
	Storage        Storage
	stateHandlers  map[models.State]HandlerFunc
	defaultHandler HandlerFunc
}

// NewRouter creates a new router
func NewRouter(storage Storage) *Router {
	return &Router{
		Storage:       storage,
		stateHandlers: make(map[models.State]HandlerFunc),
	}
}

// Message registers a handler for a specific state (similar to aiogram's router.message decorator)
func (r *Router) Register(state models.State, handler HandlerFunc) {
	r.stateHandlers[state] = handler
}

// DefaultMessage sets a default message handler for unhandled states
func (r *Router) DefaultMessage(handler HandlerFunc) {
	r.defaultHandler = handler
}

// ProcessUpdate processes an update based on the current FSM state (similar to aiogram's Dispatcher)
func (r *Router) ProcessUpdate(ctx context.Context, update *tgbotapi.Update, bot *models.Bot) error {
	if update.Message == nil {
		return nil // Only process messages for now
	}

	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	message := update.Message

	fsm := NewFSMContext(ctx, r.Storage, chatID, userID)

	state, err := fsm.Current()
	if err != nil {
		return err
	}

	handler, exists := r.stateHandlers[state]
	if !exists {
		if r.defaultHandler != nil {
			return r.defaultHandler(ctx, fsm, message, bot)
		}
		return errors.New("no handler for state")
	}

	return handler(ctx, fsm, message, bot)
}

// // Start the bot and begin processing updates
// func (r *Router) Start(bot *Bot) {
// 	ctx := context.Background()

// 	for update := range bot.UpdateChannel {
// 		if err := r.ProcessUpdate(ctx, &update, bot); err != nil {
// 			log.Printf("Error processing update: %v", err)
// 		}
// 	}
// }

// Define your states
const (
	DefaultState       models.State = ""
	StateStart         models.State = "start"
	StateAwaitingLogin models.State = "awaiting_login"
	StateAwaitingOTP   models.State = "awaiting_otp"
	StateRegistered    models.State = "registered"
)
