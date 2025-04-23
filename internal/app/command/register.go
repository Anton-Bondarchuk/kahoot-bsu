package command

import (
	"kahoot_bsu/internal/service/fsm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RegisterCommand struct {
	*CommandHandler
	fsm *fsm.FSMContext
}

func NewRegisterCommand(
	commandHandler *CommandHandler,
	fsm *fsm.FSMContext,
) *RegisterCommand {
	return &RegisterCommand{
		CommandHandler: commandHandler,
		fsm:            fsm,
	}
}

func (h *RegisterCommand) Execute(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "🔑 Пожалуйста, введите ваш login для регистрации.")

	h.bot.Telegram.Send(msg)
	h.fsm.Set(fsm.StateAwaitingLogin)
}
