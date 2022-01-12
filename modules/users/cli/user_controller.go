package cli

import (
	"api/modules/users/usecases"
	"fmt"
)

type UserController interface {
	AddUserGenesis() error
}

type UserControllerCli struct {
	userUsecase usecases.UserUsecase
}

func NewUserController(userUsecase usecases.UserUsecase) UserController {
	return &UserControllerCli{userUsecase}
}

func (controller *UserControllerCli) AddUserGenesis() error {
	userCreated, usecaseErr, serverErr := controller.userUsecase.CreateGenesisUser()
	if serverErr != nil {
		return serverErr
	}
	if usecaseErr != nil {
		fmt.Println(usecaseErr)
		return nil
	}
	fmt.Println(userCreated)
	return nil
}
