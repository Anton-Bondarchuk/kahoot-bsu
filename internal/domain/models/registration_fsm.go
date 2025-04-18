package models

import "context"

type RegistrationFSM struct {
	UserID       int64 `json:"user_id"`
	WaitLogin    bool   `json:"wait_login"`
	WaitOTP      bool   `json:"wait_otp"`
	IsRegistered bool   `json:"is_registred"`
}

type RegistrationFSMRepository interface {
	UpdateOrCreate(ctx context.Context, fsm *RegistrationFSM) error
}
