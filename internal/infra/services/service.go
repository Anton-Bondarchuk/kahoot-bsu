package services

import (
	"fmt"
	"kahoot_bsu/internal/infra/clients"
	"time"
)

type EmailService struct {
	Client EmailClientInterface
}

type EmailClientInterface interface {
	Send(login, subject string, data map[string]any) error
}

func NewEmailService(client *clients.EmailClient) *EmailService {
	return &EmailService{
		Client: client,
	}
}

func (s *EmailService) Send(login, subject, code string, expiresAt time.Time) error {
	data := map[string]any{
		"Login":     login,
		"Code":      code,
		"ExpiresIn": fmt.Sprintf("%.0f minutes", time.Until(expiresAt).Minutes()),
	}

	return s.Client.Send(login, subject, data)
}
