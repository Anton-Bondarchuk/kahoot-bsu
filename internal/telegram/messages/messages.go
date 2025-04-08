package messages

import (
	"kahoot_bsu/internal/domain/models"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)


type MessagesHandler struct {
	bot *models.Bot
}

func New(bot *models.Bot) *MessagesHandler {
	return &MessagesHandler{
		bot: bot}
}

func (h *MessagesHandler) HandleEmailRegistration(message *tgbotapi.Message) {
	email := message.Text
	pattern := `^[^@]+@bsu\.by$`
	
	matched, err := regexp.MatchString(pattern, email)
	if err != nil || !matched {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Неверный формат email.")
		h.bot.Telegram.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, email)
	h.bot.Telegram.Send(msg)
}
