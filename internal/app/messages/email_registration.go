package messages

import (
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


type HandleEmailRegistrationMessenger struct {
	*MessagesHandler
}

func (h *MessagesHandler) Execute(message *tgbotapi.Message) {
	login := message.Text
	pattern := "rct.+"
	
	matched, err := regexp.MatchString(pattern, login)
	if err != nil || !matched {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Неверный формат login.")
		h.bot.Telegram.Send(msg)
		return
	}

	// logic for add new verification code 
	msg := tgbotapi.NewMessage(message.Chat.ID, login)
	h.bot.Telegram.Send(msg)
}
