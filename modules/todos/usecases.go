package todos

import (
	usersUsecase "api/modules/users/usecases"
	"bytes"
	"errors"
)

type TodoUsecase interface {
	// TODO: Mudar parâmetros de todas funções para dto (data transfer object)
	CreateTodo(title, description string, statusTodoId, userId int64) (todo *Todo, usecaseErr error, serverErr error)
	UpdateTodo(todoID int64, title, description string, statusTodoId, userId int64) (usecaseErr error, serverErr error)
	DeleteTodo(todoID int64) (usecaseErr error, serverErr error)
	GetTodo(todoID int64) (todo *Todo, usecaseErr error, serverErr error)
	GetAllTodo() (todos []*Todo, usecaseErr error, serverErr error)

	UpdateImageTodo(dto *UpdateImageTodoDTO) (usecaseErr error, serverErr error)
	GetImageTodo(todoID int64) (image *bytes.Buffer, usecaseErr error, serverErr error)
	DeleteImageTodo(todoID int64) (usecaseErr error, serverErr error)

	CreateStatusTodo(name string, userId int64) (statusTodo *StatusTodo, usecaseErr error, serverErr error)
	UpdateStatusTodo(userId int64, statusTodoId int64, name string) (usecaseErr error, serverErr error)
	GetStatusTodo(userId, id int64) (statusTodo *StatusTodo, usecaseErr error, serverErr error)
	GetAllStatusTodo() (allStatusTodo []*StatusTodo, usecaseErr error, serverErr error)
	DeleteStatusTodo(userId, statusId int64) (usecaseErr error, serverErr error)
}

var (
	ErrTitleIsLong             = errors.New("title is too long")
	ErrDescriptionIsLong       = errors.New("description is too long")
	ErrStatusTodoNotFound      = errors.New("status todo not found")
	ErrNameStatusTodoIsSmall   = errors.New("name status is small")
	ErrStatusTodoAlreadyExists = errors.New("status já existe")
	ErrStatusTodoIdNegative    = errors.New("status todo id should to be positive")
	ErrUserIdNegative          = errors.New("user id should to be positive")
	ErrTodoNotFound            = errors.New("todo not found")
	ErrUserNotFound            = errors.New("user not found")
	ErrTodoIdIsNegative        = errors.New("todo id should be positive")
	ErrHasTodosWithStatusId    = errors.New("essa lista tem alguns Item, remove-os antes")
	ErrImageNotFound           = errors.New("image not found")
)

type DBTodoUsecase struct {
	todoRepository TodoRepository
	userRepository usersUsecase.UserUsecase
}

func NewTodoUsecase(
	todoRepository TodoRepository,
	userRepository usersUsecase.UserUsecase,
) TodoUsecase {
	return &DBTodoUsecase{todoRepository, userRepository}
}

func (usecase *DBTodoUsecase) CreateTodo(title, description string, statusTodoId, userId int64) (todo *Todo, usecaseErr error, serverErr error) {
	// TODO: mover validação para o model

	if len(title) > 255 {
		usecaseErr = ErrTitleIsLong
		return
	}
	if len(description) > 255 {
		usecaseErr = ErrDescriptionIsLong
		return
	}
	if statusTodoId <= 0 {
		usecaseErr = ErrStatusTodoIdNegative
		return
	}

	statusFound, err := usecase.todoRepository.GetStatusTodo(userId, statusTodoId)
	if err != nil {
		serverErr = err
		return
	}
	if statusFound == nil {
		usecaseErr = ErrStatusTodoNotFound
		return
	}

	todo, err = usecase.todoRepository.InsertTodo(title, description, statusTodoId)
	if err != nil {
		serverErr = err
		return
	}

	return
}

func (usecase *DBTodoUsecase) UpdateTodo(todoId int64, title, description string, statusTodoId, userId int64) (usecaseErr error, serverErr error) {

	// TODO: mover validação para o model
	if statusTodoId <= 0 {
		usecaseErr = ErrStatusTodoIdNegative
		return
	}

	if todoId <= 0 {
		usecaseErr = ErrTodoIdIsNegative
		return
	}

	if len(title) > 255 {
		usecaseErr = ErrTitleIsLong
		return
	}
	if len(description) > 255 {
		usecaseErr = ErrDescriptionIsLong
		return
	}

	statusFound, usecaseErr, serverErr := usecase.GetStatusTodo(userId, statusTodoId)
	if serverErr != nil {
		return
	}
	if usecaseErr != nil {
		return
	}
	if statusFound == nil {
		usecaseErr = ErrStatusTodoNotFound
		return
	}

	todoFound, err := usecase.todoRepository.GetTodo(todoId)
	if err != nil {
		serverErr = err
		return
	}
	if todoFound == nil {
		usecaseErr = ErrTodoNotFound
		return
	}

	err = usecase.todoRepository.UpdateTodo(todoId, title, description, statusTodoId)
	if err != nil {
		serverErr = err
		return
	}
	return
}

func (usecase *DBTodoUsecase) DeleteTodo(todoID int64) (usecaseErr error, serverErr error) {
	if todoID <= 0 {
		usecaseErr = ErrTodoIdIsNegative
		return
	}

	todoFound, usecaseErr, serverErr := usecase.GetTodo(todoID)
	if usecaseErr != nil || serverErr != nil {
		return
	}
	if todoFound == nil {
		usecaseErr = ErrTodoNotFound
		return
	}
	serverErr = usecase.todoRepository.DeleteTodo(todoID)
	return
}

func (usecase *DBTodoUsecase) GetTodo(todoID int64) (todo *Todo, usecaseErr error, serverErr error) {
	if todoID <= 0 {
		usecaseErr = ErrTodoIdIsNegative
		return
	}
	todo, serverErr = usecase.todoRepository.GetTodo(todoID)
	return
}

func (usecase *DBTodoUsecase) GetAllTodo() (todos []*Todo, usecaseErr error, serverErr error) {
	todos, serverErr = usecase.todoRepository.GetAllTodo()
	return
}

func (usecase *DBTodoUsecase) UpdateImageTodo(dto *UpdateImageTodoDTO) (usecaseErr error, serverErr error) {
	todoFound, usecaseErr, serverErr := usecase.GetTodo(dto.TodoId)
	if usecaseErr != nil || serverErr != nil {
		return
	}
	if todoFound == nil {
		usecaseErr = ErrTodoNotFound
		return
	}

	serverErr = usecase.todoRepository.UpdateImageTodo(dto.TodoId, &dto.BufferFile)
	return
}

func (usecase *DBTodoUsecase) GetImageTodo(todoId int64) (image *bytes.Buffer, usecaseErr error, serverErr error) {
	todoFound, usecaseErr, serverErr := usecase.GetTodo(todoId)
	if usecaseErr != nil || serverErr != nil {
		return
	}
	if todoFound == nil {
		usecaseErr = ErrTodoNotFound
		return
	}

	imageFound, serverErr := usecase.todoRepository.GetImageTodo(todoId)
	if serverErr != nil {
		return
	}
	if imageFound == nil {
		usecaseErr = ErrImageNotFound
		return
	}
	image = imageFound
	return
}

func (usecase *DBTodoUsecase) DeleteImageTodo(todoID int64) (usecaseErr error, serverErr error) {
	todoFound, usecaseErr, serverErr := usecase.GetTodo(todoID)
	if usecaseErr != nil || serverErr != nil {
		return
	}
	if todoFound == nil {
		usecaseErr = ErrTodoNotFound
		return
	}

	serverErr = usecase.todoRepository.UpdateImageTodo(todoID, nil)
	return
}

func (usecase *DBTodoUsecase) CreateStatusTodo(name string, userId int64) (statusTodo *StatusTodo, usecaseErr error, serverErr error) {
	if userId <= 0 {
		usecaseErr = ErrUserIdNegative
		return
	}

	if len(name) < 2 || len(name) > 255 {
		usecaseErr = ErrNameStatusTodoIsSmall
		return
	}

	statusTodoFound, err := usecase.todoRepository.GetStatusTodoByName(userId, name)
	if err != nil {
		serverErr = err
		return
	}
	if statusTodoFound != nil {
		usecaseErr = ErrStatusTodoAlreadyExists
		return
	}

	statusTodo, err = usecase.todoRepository.InsertStatusTodo(name, userId)
	if err != nil {
		serverErr = err
		return
	}
	return
}

func (usecase *DBTodoUsecase) UpdateStatusTodo(userId, statusTodoId int64, name string) (usecaseErr error, serverErr error) {
	if userId <= 0 {
		usecaseErr = ErrUserIdNegative
		return
	}

	if len(name) < 2 || len(name) > 255 {
		usecaseErr = ErrNameStatusTodoIsSmall
		return
	}

	if statusTodoId <= 0 {
		usecaseErr = ErrStatusTodoIdNegative
		return
	}

	statusTodoFound, usecaseErr, serverErr := usecase.GetStatusTodo(userId, statusTodoId)
	if usecaseErr != nil || serverErr != nil {
		return
	}
	if statusTodoFound == nil {
		usecaseErr = ErrStatusTodoNotFound
		return
	}

	statusTodoFoundByName, serverErr := usecase.todoRepository.GetStatusTodoByName(userId, name)
	if serverErr != nil {
		return
	}
	if statusTodoFoundByName != nil && statusTodoFoundByName.ID != statusTodoFound.ID {
		if statusTodoFoundByName.Name == name {
			usecaseErr = ErrStatusTodoAlreadyExists
			return
		}
	}

	serverErr = usecase.todoRepository.UpdateStatusTodo(statusTodoId, name, userId)
	return
}

func (usecase *DBTodoUsecase) GetStatusTodo(userId, statusId int64) (statusTodo *StatusTodo, usecaseErr error, serverErr error) {
	if statusId <= 0 {
		usecaseErr = ErrStatusTodoIdNegative
		return
	}
	if userId <= 0 {
		usecaseErr = ErrUserIdNegative
		return
	}

	statusTodo, err := usecase.todoRepository.GetStatusTodo(userId, statusId)
	if err != nil {
		serverErr = err
		return
	}

	return
}

func (usecase *DBTodoUsecase) GetAllStatusTodo() (allStatusTodo []*StatusTodo, usecaseErr error, serverErr error) {
	allStatusTodo, serverErr = usecase.todoRepository.GetAllStatusTodo()
	return
}

func (usecase *DBTodoUsecase) DeleteStatusTodo(userId, statusId int64) (usecaseErr error, serverErr error) {
	if statusId <= 0 {
		usecaseErr = ErrStatusTodoIdNegative
		return
	}

	if userId <= 0 {
		usecaseErr = ErrUserIdNegative
		return
	}

	statusTodoFound, usecaseErr, serverErr := usecase.GetStatusTodo(userId, statusId)
	if usecaseErr != nil || serverErr != nil {
		return
	}
	if statusTodoFound == nil {
		usecaseErr = ErrStatusTodoNotFound
		return
	}

	countTodoOnStatusTodo, serverErr := usecase.todoRepository.CountTodoByStatus(statusId)
	if usecaseErr != nil {
		return
	}
	if countTodoOnStatusTodo > 0 {
		usecaseErr = ErrHasTodosWithStatusId
		return
	}

	serverErr = usecase.todoRepository.DeleteStatusTodo(userId, statusId)
	return
}
