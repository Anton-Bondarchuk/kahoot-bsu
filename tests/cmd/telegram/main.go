package main

import (
	"context"
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

	ctx := context.Background()

	telegram.Start(ctx, app)
}
