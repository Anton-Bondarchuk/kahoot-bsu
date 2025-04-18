package messages

import (
	"kahoot_bsu/internal/domain/models"
)

type MessagesHandler struct {
	bot *models.Bot
}

func New(bot *models.Bot) *MessagesHandler {
	return &MessagesHandler{
		bot: bot}
}
