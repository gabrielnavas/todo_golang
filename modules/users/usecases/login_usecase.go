package usecases

import (
	"api/modules/users/dto"
	"api/modules/users/models"
	"errors"
)

var (
	ErrCredencialsWrong = errors.New("nome de usuário ou senha inválidos")
)

type LoginUsecase interface {
	Login(username string, password string) (loginResponse dto.LoginResponse, usecaseError, serverError error)
	Logout(token string) (tokenInvalid string, usecaseError, serverError error)
}

type TokenLoginUsecase struct {
	userRepository models.UserRepository
	hashPassword   HashPassword
	tokenManager   TokenManager
}

func NewTokenLoginUsecase(
	userRepository models.UserRepository,
	hashPassword HashPassword,
	tokenManager TokenManager,
) LoginUsecase {
	return &TokenLoginUsecase{userRepository, hashPassword, tokenManager}
}

func (usecase *TokenLoginUsecase) Login(username string, password string) (loginResponse dto.LoginResponse, usecaseError, serverError error) {
	userFound, serverError := usecase.userRepository.GetUserByUsername(username)
	if serverError != nil {
		return
	}
	if userFound == nil {
		usecaseError = ErrCredencialsWrong
		return
	}

	passwordEquals, serverError := usecase.hashPassword.Verify(password, userFound.Password)
	if serverError != nil {
		return
	}
	if !passwordEquals {
		usecaseError = ErrCredencialsWrong
		return
	}

	token, serverError := usecase.tokenManager.CreateToken(
		userFound.ID,
		userFound.LevelAccess,
		models.DurationTimeDefault,
	)
	if serverError != nil {
		return
	}
	loginResponse = dto.LoginResponse{
		Token: token,
		User:  *userFound,
	}

	return
}

func (usecase *TokenLoginUsecase) Logout(token string) (tokenInvalid string, usecaseError, serverError error) {
	return "", nil, nil
}
