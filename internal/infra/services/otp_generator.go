package services

import (
	"crypto/rand"
	contract "kahoot_bsu/internal/ports/service"
)

// otp - one time password
type verificationOTPGenerator struct {
	lenght int
}

func NewVerificationOTPGenerator(lenght int) contract.VerificationCodeGenerator {
	return &verificationOTPGenerator{
		lenght: lenght,
	}
}

func (v *verificationOTPGenerator) Generate() (string, error) {
	const otpChars = "1234567890"

	buffer := make([]byte, v.lenght)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := range v.lenght {
		buffer[i] = otpChars[int(buffer[i]) % otpCharsLength]
	}

	return string(buffer), nil
}
