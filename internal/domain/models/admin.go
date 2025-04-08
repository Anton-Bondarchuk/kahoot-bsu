package models


type Admin struct {
	TelegramID int64        `json:"telegram_id"`
	Name       string       `json:"name"`
	Email      string       `json:"email"`
}
