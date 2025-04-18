package command

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type UnknownCommand struct {
	*CommandHandler
}

func (h *UnknownCommand) Execute(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "🔑 Пожалуйста, введите ваш email для регистрации.")
	h.bot.Telegram.Send(msg)
}
