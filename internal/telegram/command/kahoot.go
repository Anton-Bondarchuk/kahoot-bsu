package command

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type KahootComand struct {
	*CommandHandler
}

func (h *KahootComand) Execute(message *tgbotapi.Message) {
	kahootMsgText := "Нажмите на кнопку ниже, чтобы запустить приложение"
	kbRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonWebApp("Kahoot!", tgbotapi.WebAppInfo{URL: h.WebAppUrl}),
	)

	kb := tgbotapi.NewInlineKeyboardMarkup(kbRow)

	msg := tgbotapi.NewMessage(message.Chat.ID, kahootMsgText)
	msg.ReplyMarkup = kb

	h.bot.Telegram.Send(msg)
}
