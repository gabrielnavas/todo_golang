package todos

import (
	"bytes"
	"errors"
	"strings"
)

type TodoUsecase interface {
	CreateTodo(title, description string, statusTodoId int64) (todo *Todo, usecaseErr error, serverErr error)
	UpdateTodo(todoID int64, title, description string, statusTodoId int64) (usecaseErr error, serverErr error)
	DeleteTodo(todoID int64) (usecaseErr error, serverErr error)
	GetTodo(todoID int64) (todo *Todo, usecaseErr error, serverErr error)
	GetAllTodo() (todos []*Todo, usecaseErr error, serverErr error)

	UpdateImageTodo(dto *UpdateImageTodoDTO) (usecaseErr error, serverErr error)
	GetImageTodo(todoID int64) (image *bytes.Buffer, usecaseErr error, serverErr error)
	// DeleteImageTodo(todoID int64) (usecaseErr error, serverErr error)

	CreateStatusTodo(name string) (statusTodo *StatusTodo, usecaseErr error, serverErr error)
	UpdateStatusTodo(statusTodoId int64, name string) (usecaseErr error, serverErr error)
	GetStatusTodo(id int64) (statusTodo *StatusTodo, usecaseErr error, serverErr error)
	GetAllStatusTodo() (allStatusTodo []*StatusTodo, usecaseErr error, serverErr error)
	DeleteStatusTodo(id int64) (usecaseErr error, serverErr error)
}

var (
	ErrTitleIsLong             = errors.New("title is too long")
	ErrDescriptionIsLong       = errors.New("description is too long")
	ErrStatusTodoNotFound      = errors.New("status todo not found")
	ErrNameStatusTodoIsSmall   = errors.New("name status is small")
	ErrStatusTodoAlreadyExists = errors.New("status todo already exists")
	ErrStatusTodoIdNegative    = errors.New("status todo id should to be positive")
	ErrTodoNotFound            = errors.New("todo not found")
	ErrTodoIdIsNegative        = errors.New("todo id should be positive")
	ErrHasTodosWithStatusId    = errors.New("has todos with this status id")
	ErrImageNotFound           = errors.New("image not found")
)

type DBTodoUsecase struct {
	todoRepository TodoRepository
}

func NewTodoUsecase(todoRepository TodoRepository) TodoUsecase {
	return &DBTodoUsecase{todoRepository}
}

func (usecase *DBTodoUsecase) CreateTodo(title, description string, statusTodoId int64) (todo *Todo, usecaseErr error, serverErr error) {
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

	statusFound, err := usecase.todoRepository.GetStatusTodo(statusTodoId)
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

func (usecase *DBTodoUsecase) UpdateTodo(todoID int64, title, description string, statusTodoId int64) (usecaseErr error, serverErr error) {
	if statusTodoId <= 0 {
		usecaseErr = ErrStatusTodoIdNegative
		return
	}

	if todoID <= 0 {
		usecaseErr = ErrTodoIdIsNegative
		return
	}

	statusFound, usecaseErr, serverErr := usecase.GetStatusTodo(statusTodoId)
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

	todoFound, err := usecase.todoRepository.GetTodo(todoID)
	if err != nil {
		serverErr = err
		return
	}
	if todoFound == nil {
		usecaseErr = ErrTodoNotFound
		return
	}

	err = usecase.todoRepository.UpdateTodo(todoID, title, description, statusTodoId)
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

func (usecase *DBTodoUsecase) GetImageTodo(todoID int64) (image *bytes.Buffer, usecaseErr error, serverErr error) {
	imageFound, serverErr := usecase.todoRepository.GetImageTodo(todoID)
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

func (usecase *DBTodoUsecase) CreateStatusTodo(name string) (statusTodo *StatusTodo, usecaseErr error, serverErr error) {
	if len(name) < 2 || len(name) > 255 {
		usecaseErr = ErrNameStatusTodoIsSmall
		return
	}

	nameLower := strings.ToLower(name)

	statusTodoFound, err := usecase.todoRepository.GetStatusTodoByName(nameLower)
	if err != nil {
		serverErr = err
		return
	}
	if statusTodoFound != nil {
		usecaseErr = ErrStatusTodoAlreadyExists
		return
	}

	statusTodo, err = usecase.todoRepository.InsertStatusTodo(nameLower)
	if err != nil {
		serverErr = err
		return
	}
	return
}

func (usecase *DBTodoUsecase) UpdateStatusTodo(statusTodoId int64, name string) (usecaseErr error, serverErr error) {
	if statusTodoId <= 0 {
		usecaseErr = ErrStatusTodoIdNegative
		return
	}

	statusTodoFound, usecaseErr, serverErr := usecase.GetStatusTodo(statusTodoId)
	if usecaseErr != nil || serverErr != nil {
		return
	}
	if statusTodoFound == nil {
		usecaseErr = ErrStatusTodoNotFound
		return
	}

	// tratar nomes iguais
	nameLower := strings.ToLower(name)
	statusTodoFoundByName, serverErr := usecase.todoRepository.GetStatusTodoByName(nameLower)
	if serverErr != nil {
		return
	}
	if statusTodoFoundByName != nil && statusTodoFoundByName.ID != statusTodoFound.ID {
		if statusTodoFoundByName.Name == nameLower {
			usecaseErr = ErrStatusTodoAlreadyExists
			return
		}
	}

	serverErr = usecase.todoRepository.UpdateStatusTodo(statusTodoId, nameLower)
	return
}

func (usecase *DBTodoUsecase) GetStatusTodo(id int64) (statusTodo *StatusTodo, usecaseErr error, serverErr error) {
	if id <= 0 {
		usecaseErr = ErrStatusTodoIdNegative
		return
	}

	statusTodo, err := usecase.todoRepository.GetStatusTodo(id)
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

func (usecase *DBTodoUsecase) DeleteStatusTodo(id int64) (usecaseErr error, serverErr error) {
	if id <= 0 {
		usecaseErr = ErrStatusTodoIdNegative
		return
	}

	statusTodoFound, usecaseErr, serverErr := usecase.GetStatusTodo(id)
	if usecaseErr != nil || serverErr != nil {
		return
	}
	if statusTodoFound == nil {
		usecaseErr = ErrStatusTodoNotFound
		return
	}

	countTodoOnStatusTodo, serverErr := usecase.todoRepository.CountTodoByStatus(id)
	if usecaseErr != nil {
		return
	}
	if countTodoOnStatusTodo > 0 {
		usecaseErr = ErrHasTodosWithStatusId
		return
	}

	serverErr = usecase.todoRepository.DeleteStatusTodo(id)
	return
}
