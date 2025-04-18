package command

import (
	"context"
	"kahoot_bsu/internal/domain/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RegisterCommand struct {
	*CommandHandler
	FSM models.RegistrationFSMRepository 
}

func NewRegisterCommand(
	commandHandler *CommandHandler,
	fsm models.RegistrationFSMRepository, 
) * RegisterCommand {
	return &RegisterCommand{
		CommandHandler: commandHandler,
		FSM: fsm,
	}
}

func (h *RegisterCommand) Execute(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "🔑 Пожалуйста, введите ваш login для регистрации.")
	fsmModel := &models.RegistrationFSM{
		UserID: h.bot.Telegram.Self.ID,
		WaitLogin: true,
	}
	h.FSM.UpdateOrCreate(context.Background(), fsmModel)
	h.bot.Telegram.Send(msg)
}
