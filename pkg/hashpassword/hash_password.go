package hashpassword

import "golang.org/x/crypto/bcrypt"

type HashPassword interface {
	Hash(passwordPlain string) (string, error)
	Verify(passwordPlain, passwordHashed string) (bool, error)
}

type HashPasswordBcrypt struct{}

func NewHashPassword() HashPassword {
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
