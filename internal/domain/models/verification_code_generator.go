package models

type VerificationCodeGenerator interface {
	Generate() (string, error)
}