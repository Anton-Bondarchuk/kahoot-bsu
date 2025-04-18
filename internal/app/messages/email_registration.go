package messages

import (
	"context"
	"fmt"
	"kahoot_bsu/internal/auth"
	"kahoot_bsu/internal/domain/models"
	"kahoot_bsu/internal/service/email"
	"log"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO: refactor handler
type EmailRegistrationHandler struct {
	handler      *MessagesHandler
	emailService EmailSender
	codeGen      models.VerificationCodeGenerator
	userRepo     models.UserRepository
	fSMService   FSMService
}

type EmailSender interface {
	Send(login, subject, code string, expiresAt time.Time) error
}

type FSMService interface {
	SetState(ctx context.Context, userID int64, code string, state *models.RegistrationFSM) error
}

// TODO: review and fix interfaces
func NewEmailRegistrationHandler(
	messagesHandler *MessagesHandler,
	emailService *email.EmailService,
	codeGen      models.VerificationCodeGenerator,
	userRepo models.UserRepository,
	fsmSerive FSMService,
) *EmailRegistrationHandler {
	return &EmailRegistrationHandler{
		handler:      messagesHandler,
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (h *EmailRegistrationHandler) Execute(message *tgbotapi.Message) {
	login := message.Text
	pattern := "rct.+|bio.+"

	matched, err := regexp.MatchString(pattern, login)
	if err != nil || !matched {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Неверный формат login.")
		h.handler.bot.Telegram.Send(msg)
		return
	}

	userId := h.handler.bot.Telegram.Self.ID
	user := models.User{
		ID:    userId,
		Login: login,
		Role:  int64(auth.RoleBlocked),
	}
	// TODO: use messages context?

	err = h.userRepo.UpdateOrCreate(context.Background(), &user)
	code, err := h.codeGen.Generate()
	if err != nil {
		log.Printf("%e", err)
	}
	
	expiresAt := time.Now().Add(30 * time.Minute)

	if err := h.emailService.Send(login[4:], "Your Verification Code", code, expiresAt); err != nil {
		log.Printf("Failed to sentd verification email: %v", err)
	} else {
		log.Printf("Verification email sent to %s", login)
	}

	userID := h.handler.bot.Telegram.Self.ID

	// TODO: Add request to the database and update only nessessary fields
	err = h.fSMService.SetState(context.Background(), userID, code, &models.RegistrationFSM{
		UserID: userID,
		WaitLogin: true,
		WaitOTP: true,
	})

	// TODO: mail addres to constant
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Вам на почту https://webmail.bsu.by/owa/#path=/mail был отправлен код, введите его, чтобы подвердить, что вы студент БГУ"))
	h.handler.bot.Telegram.Send(msg)
}
