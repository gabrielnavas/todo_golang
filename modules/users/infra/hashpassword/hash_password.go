package hashpassword

import (
	"api/modules/users/usecases"

	"golang.org/x/crypto/bcrypt"
)

type HashPasswordBcrypt struct{}

func NewHashPassword() usecases.HashPassword {
	return &HashPasswordBcrypt{}
}

func (hashp *HashPasswordBcrypt) Hash(passwordPlain string) (string, error) {
	var passwordHashed []byte
	var err error

	passwordHashed, err = bcrypt.GenerateFromPassword([]byte(passwordPlain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordHashed), nil
}

func (hashp *HashPasswordBcrypt) Verify(passwordPlain, passwordHashed string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(passwordPlain))
	if err == nil {
		return true, nil
	}
	return false, nil
}
