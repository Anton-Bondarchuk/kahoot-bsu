package telegram_test

// import (
// 	"kahoot_bsu/internal/domain/models"
// 	"kahoot_bsu/internal/app/command"
// 	"kahoot_bsu/internal/app/messages"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// func Start(b *models.Bot) {
// 	for update := range b.UpdateChannel {
// 		if update.Message == nil {
// 			continue
// 		}

// 		if update.Message.IsCommand() {
// 			handleCommand(b, update.Message)
// 		} else {
// 			handleMessages(b, update.Message)
// 		}
// 	}
// }

// type CommandInterface interface {
// 	Execute(message *tgbotapi.Message)
// }

// func handleCommand(b *models.Bot, message *tgbotapi.Message) {
// 	comandHandler := command.New(b, "https://af09-185-53-133-77.ngrok-free.app/")

// 	commandStrategy := map[string]CommandInterface{
// 		"start": &command.StartCommand{CommandHandler: comandHandler},
// 		"register":   &command.RegisterCommand{CommandHandler: comandHandler},
// 		"kahoot":   &command.KahootComand{CommandHandler: comandHandler},
// 		"help": &command.HelpCommand{CommandHandler: comandHandler},
// 		"unknown": &command.UnknownCommand{CommandHandler: comandHandler},
// 	}

// 	_, ok := commandStrategy[message.Command()]

// 	if (!ok) {
// 		commandStrategy["unknown"].Execute(message)
// 		return
// 	}

// 	commandStrategy[message.Command()].Execute(message)
// }

// func handleMessages(b *models.Bot, message *tgbotapi.Message) {
// 	messHandler := messages.New(b)

// 	emailReg := &messages.HandleEmailRegistrationMessenger{MessagesHandler: messHandler}
// 	emailReg.Execute(message)
// }