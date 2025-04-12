package command

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type RegisterCommand struct {
	*CommandHandler
}
func (h *RegisterCommand) Execute(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "🔑 Пожалуйста, введите ваш login для регистрации.")
	h.bot.Telegram.Send(msg)
}
