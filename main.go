package main

import (
	"log"
	"time"

	"kahoot_bsu/internal/service/email"
)

func main() {
	// Create an email service with configuration
	emailService := email.NewEmailService(email.Config{
		Host:       "smtp.bsu.by",     // e.g., "smtp.gmail.com"
		Port:       587,                        // Common TLS port
		Username:   "rct.bondarchAS@bsu.by", // Your email
		Password:   "", // Your password or app password
		FromEmail:  "rct.bondarchAS@bsu.by",
		FromName:   "BSU Quiz Platform",
		Domain:     "bsu.by",
		Prefix:     "rct.",
		TemplateDir: "templates/email",
		Debug:      false, // Set to true for development
	})

	login := "bondarchAS"
	code := "123456"
	expiresAt := time.Now().Add(30 * time.Minute)

	if err := emailService.SendVerificationEmail(login, code, expiresAt); err != nil {
		log.Printf("Failed to sentd verification email: %v", err)
	} else {
		log.Printf("Verification email sent to %s", emailService.FormatBSUEmail(login))
	}
}