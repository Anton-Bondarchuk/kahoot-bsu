package service

type VerificationCodeGenerator interface {
	Generate() (string, error)
}