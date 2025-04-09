package menu

import (
	"kahoot_bsu/internal/domain/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SetMenu(b *models.Bot, commands []tgbotapi.BotCommand) (*tgbotapi.APIResponse, error) {
	cfg := tgbotapi.NewSetMyCommands(commands...)

	return b.Telegram.Request(cfg)
}