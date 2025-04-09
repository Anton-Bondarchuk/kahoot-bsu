package main

import (
	"kahoot_bsu/internal/bot"
	"kahoot_bsu/internal/config"
	"log"
)

func main() {
	cfg := config.MustLoad()

	telegramBot, err := bot.New(cfg.BotConfig)
	if err != nil {
		log.Fatalf("Exception on the start the bot: %v", err)
	}

	bot.Start(telegramBot)
}