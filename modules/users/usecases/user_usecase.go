package usecases

import (
	"bytes"
	"errors"
	"time"

	"api/modules/users/dto"
	"api/modules/users/models"
)

var (
	ErrUserAlreadyExistsWithEmail            = errors.New("usuário já existe com esse email")
	ErrUserAlreadyExistsWithUsername         = errors.New("usuário já existe com esse nome de usuário")
	ErrPasswordNotEqualsPasswordConfirmation = errors.New("password is not equals password confirmation")
	ErrUserNotFound                          = errors.New("user not found")
	ErrPhotoNotFound                         = errors.New("photo not found")
	ErrOldPasswordWrong                      = errors.New("senha antiga inválida")
	ErrAllreadyHaveUsers                     = errors.New("allready have users")
)

type UserUsecase interface {
	CreateGenesisUser() (userCreated *models.User, usecaseError, serverError error)
	CreateUser(name, username, password, passwordConfirmation, email string) (userCreated *models.User, usecaseError, serverError error)
	UpdateUser(id int64, name, username, password, email string, levelAccess models.LevelAccess) (usecaseError, serverError error)
	DeleteUser(id int64) (usecaseError, serverError error)
	GetUser(id int64) (userFound *models.User, usecaseError, serverError error)
	GetAllUser() (userFound []*models.User, usecaseError, serverError error)
	ChangePassword(userId int64, oldPassword, newPassword, newPasswordConfirmation string) (usecaseError, serverError error)

	UpdatePhotoUser(photo *dto.UpdatePhotoUserDTO) (usecaseError, serverError error)
	DeletePhotoUser(id int64) (usecaseError, serverError error)
	GetPhotoUser(id int64) (photo *bytes.Buffer, usecaseError, serverError error)
}

type DBUserUsecase struct {
	userRepository models.UserRepository
	hashPassword   HashPassword
}

func NewUserUsecase(userRepository models.UserRepository, hashPassword HashPassword) UserUsecase {
	return &DBUserUsecase{userRepository, hashPassword}
}

func (usecase *DBUserUsecase) CreateGenesisUser() (userCreated *models.User, usecaseError, serverError error) {
	userCreated, usecaseError = models.NewUser(
		0,
		"adm",
		"adm",
		"123456",
		"adm@email.com",
		models.AdminLevelAccess,
		time.Now(),
		time.Now(),
		bytes.Buffer{},
	)
	if usecaseError != nil {
		return
	}

	userCount, serverError := usecase.userRepository.CountUser()
	if serverError != nil {
		return
	}

	if userCount > 0 {
		usecaseError = ErrAllreadyHaveUsers
		return
	}

	passwordHashed, serverError := usecase.hashPassword.Hash(userCreated.Password)
	if serverError != nil {
		return
	}

	userCreated, serverError = usecase.userRepository.InsertUser(
		userCreated.Name,
		userCreated.Username,
		passwordHashed,
		userCreated.Email,
		userCreated.LevelAccess,
	)
	if serverError != nil {
		return
	}

	return userCreated, nil, nil
}

func (usecase *DBUserUsecase) CreateUser(name, username, password, passwordConfirmation, email string) (userCreated *models.User, usecaseError, serverError error) {
	if password != passwordConfirmation {
		usecaseError = ErrPasswordNotEqualsPasswordConfirmation
		return
	}
	user, usecaseError := models.NewUser(
		0,
		name,
		username,
		password,
		email,
		models.BasicLevelAccess,
		time.Now(),
		time.Now(),
		bytes.Buffer{},
	)
	if usecaseError != nil {
		return
	}

	userFoundByEmail, serverError := usecase.userRepository.GetUserByEmail(email)
	if serverError != nil {
		return
	}
	if userFoundByEmail != nil {
		usecaseError = ErrUserAlreadyExistsWithEmail
		return
	}

	userFoundByUsername, serverError := usecase.userRepository.GetUserByUsername(username)
	if serverError != nil {
		return
	}
	if userFoundByUsername != nil {
		usecaseError = ErrUserAlreadyExistsWithUsername
		return
	}

	passwordHashed, serverError := usecase.hashPassword.Hash(password)
	if serverError != nil {
		return
	}

	userCreated, serverError = usecase.userRepository.InsertUser(name, username, passwordHashed, email, user.LevelAccess)
	if serverError != nil {
		return
	}

	return userCreated, nil, nil
}

func (usecase *DBUserUsecase) UpdateUser(id int64, name, username, password, email string, levelAccess models.LevelAccess) (usecaseError, serverError error) {
	user, usecaseError := models.NewUser(
		0,
		name,
		username,
		password,
		email,
		levelAccess,
		time.Now(),
		time.Now(),
		bytes.Buffer{},
	)
	if usecaseError != nil {
		return usecaseError, nil
	}

	userFoundByEmail, serverError := usecase.userRepository.GetUserByEmail(email)
	if serverError != nil {
		return nil, serverError
	}
	if userFoundByEmail != nil && email != userFoundByEmail.Email {
		return ErrUserAlreadyExistsWithEmail, nil
	}

	userFoundByUsername, serverError := usecase.userRepository.GetUserByUsername(username)
	if serverError != nil {
		return nil, serverError
	}
	if userFoundByUsername != nil && username != userFoundByEmail.Username {
		return ErrUserAlreadyExistsWithUsername, nil
	}

	passwordHashed, serverError := usecase.hashPassword.Hash(password)
	if serverError != nil {
		return nil, serverError
	}

	usecase.userRepository.UpdateUser(id, name, username, passwordHashed, email, user.LevelAccess)

	return nil, nil
}

func (usecase *DBUserUsecase) DeleteUser(id int64) (usecaseError, serverError error) {
	userFoundById, serverError := usecase.userRepository.GetUser(id)
	if serverError != nil {
		return nil, serverError
	}
	if userFoundById == nil {
		usecaseError = ErrUserNotFound
		return usecaseError, nil
	}

	usecase.userRepository.DeleteUser(id)

	return nil, nil
}

func (usecase *DBUserUsecase) GetUser(id int64) (userFound *models.User, usecaseError, serverError error) {
	userFoundById, serverError := usecase.userRepository.GetUser(id)
	if serverError != nil {
		return nil, nil, serverError
	}
	if userFoundById == nil {
		usecaseError = ErrUserNotFound
		return nil, usecaseError, nil
	}

	return userFoundById, nil, nil
}

func (usecase *DBUserUsecase) GetAllUser() (userFound []*models.User, usecaseError, serverError error) {
	allUsers, serverError := usecase.userRepository.GetAllUser()
	if serverError != nil {
		return nil, nil, serverError
	}
	return allUsers, nil, nil
}

func (usecase *DBUserUsecase) ChangePassword(userId int64, oldPassword, newPassword, newPasswordConfirmation string) (usecaseError, serverError error) {
	userFound, usecaseError, serverError := usecase.GetUser(userId)
	if usecaseError != nil || serverError != nil {
		return
	}
	if userFound == nil {
		usecaseError = ErrUserNotFound
		return
	}

	passwordEquals, serverError := usecase.hashPassword.Verify(oldPassword, userFound.Password)
	if serverError != nil {
		return
	}
	if !passwordEquals {
		usecaseError = ErrOldPasswordWrong
		return
	}

	if newPassword != newPasswordConfirmation {
		usecaseError = ErrPasswordNotEqualsPasswordConfirmation
		return
	}

	passwordHashed, serverError := usecase.hashPassword.Hash(newPassword)
	if serverError != nil {
		return nil, serverError
	}

	serverError = usecase.userRepository.UpdateUser(
		userFound.ID,
		userFound.Name,
		userFound.Username,
		passwordHashed,
		userFound.Email,
		userFound.LevelAccess,
	)
	if serverError != nil {
		return
	}
	return
}

func (usecase *DBUserUsecase) UpdatePhotoUser(photoDto *dto.UpdatePhotoUserDTO) (usecaseError, serverError error) {
	userFound, usecaseError, serverError := usecase.GetUser(photoDto.UserId)
	if usecaseError != nil || serverError != nil {
		return
	}
	if userFound == nil {
		usecaseError = ErrUserNotFound
		return
	}

	serverError = usecase.userRepository.UpdatePhotoUser(photoDto.UserId, &photoDto.BufferFile)
	return
}

func (usecase *DBUserUsecase) DeletePhotoUser(id int64) (usecaseError, serverError error) {
	userFound, usecaseError, serverError := usecase.GetUser(id)
	if usecaseError != nil || serverError != nil {
		return
	}
	if userFound == nil {
		usecaseError = ErrUserNotFound
		return
	}

	serverError = usecase.userRepository.UpdatePhotoUser(id, nil)
	return
}

func (usecase *DBUserUsecase) GetPhotoUser(id int64) (photo *bytes.Buffer, usecaseError, serverError error) {
	userFound, usecaseError, serverError := usecase.GetUser(id)
	if usecaseError != nil || serverError != nil {
		return
	}
	if userFound == nil {
		usecaseError = ErrUserNotFound
		return
	}

	photo, serverError = usecase.userRepository.GetPhotoUser(id)
	if serverError != nil {
		return
	}
	if photo == nil {
		usecaseError = ErrPhotoNotFound
		return
	}
	return
}
