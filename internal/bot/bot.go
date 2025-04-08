package bot

import (
	"kahoot_bsu/internal/config"
	"kahoot_bsu/internal/domain/models"
	"kahoot_bsu/internal/telegram/command"
	"kahoot_bsu/internal/telegram/messages"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Bot - структура для хранения бота и базы данных

// NewBot - функция для создания нового бота
func New(cfg config.BotConfig) (*models.Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	// TODO: add With and functional option pattern
	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Timeout

	updates, err := botAPI.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}

	return &models.Bot{
		Telegram:      botAPI,
		// DB:            db,
		UpdateChannel: updates, // Инициализируем канал обновлений
	}, nil
}

func Start(b *models.Bot) {
	messHandler := messages.New(b)
	for update := range b.UpdateChannel {
		if update.Message != nil {
			if update.Message.IsCommand() {
				handleCommand(b, update.Message)
			} else {
				messHandler.HandleEmailRegistration(update.Message)
			}
		}
	}
}

func handleCommand(b *models.Bot, message *tgbotapi.Message) {
	comandHandler := command.New(b)
	
	switch message.Command() {
	case "start":
		comandHandler.Start(message)
	case "register":
		comandHandler.Register(message)

	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда. Используйте /start или /help.")
		b.Telegram.Send(msg)
	}
}