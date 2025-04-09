package messages

import (
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


type HandleEmailRegistrationMessenger struct {
	*MessagesHandler
}

func (h *MessagesHandler) Execute(message *tgbotapi.Message) {
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
