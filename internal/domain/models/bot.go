package models 

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	Telegram      *tgbotapi.BotAPI
	UpdateChannel tgbotapi.UpdatesChannel
}