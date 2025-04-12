package telegram_test

// import (
// 	"kahoot_bsu/internal/config"
// 	"kahoot_bsu/internal/domain/models"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// type TelegramPort struct {
// 	// TODO: add some fields 
// }

// func New(cfg config.BotConfig) (*models.Bot, error) {
// 	botAPI, err := tgbotapi.NewBotAPI(cfg.Token)
// 	if err != nil {
// 		return nil, err
// 	}

// 	botAPI.Debug = cfg.Debug

// 	// TODO: add With and functional option pattern
// 	u := tgbotapi.NewUpdate(0)
// 	u.Timeout = cfg.Timeout

// 	updates := botAPI.GetUpdatesChan(u)
	
// 	return &models.Bot{
// 		Telegram:      botAPI,
// 		UpdateChannel: updates}, nil
// }


