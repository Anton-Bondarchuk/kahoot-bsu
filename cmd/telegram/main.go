package main

import (
	"kahoot_bsu/internal/service/telegram"
)

func main() {
	app, closeFunc := telegram.NewAppTelegram()
	defer func() {
		err := closeFunc()
		if err != nil {
			panic(err)
		}
	}()

	telegram.Start(app)
}
