package command

import (
	"kahoot_bsu/internal/domain/models"
) 

type CommandHandler struct {
	bot *models.Bot
	WebAppUrl string
}

func New(bot *models.Bot, WebAppUrl string) *CommandHandler {
	return &CommandHandler{
		bot: bot, 
		WebAppUrl: WebAppUrl}
}
