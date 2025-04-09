package command

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type HelpCommand struct {
	*CommandHandler
}

func (c *HelpCommand) Execute(message *tgbotapi.Message) {
	welcomeText := "👋 Добро пожаловать! Пожалуйста, зарегистрируйтесь, отправив команду /register."
	msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
	c.bot.Telegram.Send(msg)
}
