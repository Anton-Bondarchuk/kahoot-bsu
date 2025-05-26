package ports

type VerificationCodeGenerator interface {
	Generate() (string, error)
}
