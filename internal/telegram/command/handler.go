package command

import (
	"kahoot_bsu/internal/domain/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

) 

type CommandHandler struct {
	bot *models.Bot
}

func New(bot *models.Bot) *CommandHandler {
	return &CommandHandler{
		bot: bot}
}

func (h *CommandHandler) Start(message *tgbotapi.Message) {
	welcomeText := "👋 Добро пожаловать! Пожалуйста, зарегистрируйтесь, отправив команду /register."
	msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
	h.bot.Telegram.Send(msg)
}

func (h *CommandHandler) Register(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "🔑 Пожалуйста, введите ваш email для регистрации.")
	h.bot.Telegram.Send(msg)
}

func (h *CommandHandler) Kahoot(message *tgbotapi.Message) {
	// TODO: add logic for validate user permission 
	if (true) {

	}
}


