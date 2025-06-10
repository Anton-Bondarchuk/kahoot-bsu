package telegram

import (
	"context"
	"errors"
	"kahoot_bsu/internal/auth"
	"kahoot_bsu/internal/domain/models"
	"kahoot_bsu/internal/infra/services"
	ports "kahoot_bsu/internal/ports"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type fSMHandler struct {
	emailService *services.EmailService
	userRepo     ports.UserRepository
	otpGenerator ports.VerificationCodeGenerator
}

func NewFSMHandler(
	emailService *services.EmailService,
	userRepo ports.UserRepository,
	otpGenerator ports.VerificationCodeGenerator,
) *fSMHandler {
	return &fSMHandler{
		emailService: emailService,
		userRepo:     userRepo,
		otpGenerator: otpGenerator,
	}
}

func (h *fSMHandler) HandleLogin(ctx context.Context, fsm *fsmSrv.FSMContext, message *tgbotapi.Message, bot *models.Bot) error {
	login := message.Text

	// Store login in FSM data
	if err := fsm.SetData("login", login); err != nil {
		panic(err)
	}

	otp, err := h.otpGenerator.Generate()
	if err != nil {
		panic(err)
	}

	err = fsm.SetData("code", otp)

	if err != nil {
		panic(err)
	}

	expiresAt := time.Now().Add(30 * time.Minute)

	if err := h.emailService.Send(login[4:], "Your Verification Code", otp, expiresAt); err != nil {
		log.Printf("Failed to sentd verification email: %v", err)
	} else {
		log.Printf("Verification email sent to %s", login)
	}

	// msg := tgbotapi.NewMessage(message.Chat.ID, "Спасибо. На ваш [email](https://webmail.bsu.by/owa/#path=/mail) был выслан проверочный код. \n Пожалуйста, введите его:")
	// msg.ParseMode = "MarkdownV2"
	msg := tgbotapi.NewMessage(message.Chat.ID, "Спасибо! На ваш <a href=\"https://webmail.bsu.by/owa/#path=/mail\">email</a> был выслан проверочный код. \n Пожалуйста, введите его:")
	msg.ParseMode = "HTML"
	_, err = bot.Telegram.Send(msg)
	if err != nil {
		panic(err)
	}

	return fsm.Set(fsmSrv.StateAwaitingOTP)
}

func (h *fSMHandler) HandleOTP(ctx context.Context, fsm *fsmSrv.FSMContext, message *tgbotapi.Message, bot *models.Bot) error {
	inputOTP := message.Text
	fsmOTP, err := fsm.GetData("code")
	if err != nil {
		panic(err)
	}

	loginInterface, err := fsm.GetData("login")
	if err != nil {
		panic(err)
	}

	login, ok := loginInterface.(string)
	if !ok {
		return errors.New("invalid login data format")
	}

	if len(inputOTP) != 6 || inputOTP != fsmOTP {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неверный проверочный код. Введите проверчный код:")
		_, err := bot.Telegram.Send(msg)
		panic(err)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Регастрация завершена! Добро пожаловать, "+login+"!" + "\nВы можете перейти к сервису прохождения викторин кликнув на /quiz")
	_, err = bot.Telegram.Send(msg)
	if err != nil {
		panic(err)
	}

	return fsm.Set()
}

func (h *fSMHandler) HandleRegistered(ctx context.Context, fsm *fsmSrv.FSMContext, message *tgbotapi.Message, bot *models.Bot) error {
	loginInterface, err := fsm.GetData("login")
	if err != nil {
		return err
	}

	login, ok := loginInterface.(string)
	if !ok {
		return errors.New("invalid login data format")
	}

	user := &models.User{
		ID:    bot.Telegram.Self.ID,
		Login: login,
		Role:  int64(auth.RoleUser),
	}

	err = h.userRepo.UpdateOrCreate(ctx, user)
	if err != nil {
		panic(err)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Привет "+login+"! Вы уже зарегистрированы.")
	_, err = bot.Telegram.Send(msg)
	return err
}

func (h *fSMHandler) HandleDefault(ctx context.Context, fsm *fsmSrv.FSMContext, message *tgbotapi.Message, bot *models.Bot) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "I'm not sure how to respond. Try using the /start command.")
	_, err := bot.Telegram.Send(msg)
	return err
}
