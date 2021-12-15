package todos

import (
	"errors"
	"strings"
)

type TodoUsecase interface {
	CreateTodo(title, description string, statusID int64) (todo *Todo, usecaseErr error, serverErr error)
	// UpdateTodo(todoID int64, title, description string) (usecaseErr error, serverErr error)
	// DeleteTodo(todoID int64) (usecaseErr error, serverErr error)
	// GetTodo(todoID int64) (todo *Todo, usecaseErr error, serverErr error)
	GetAllTodo() (todos []*Todo, usecaseErr error, serverErr error)

	// GetImageTodo(todoID int64) (image *bytes.Buffer, serverErr error)
	// UpdateImageTodo(todoID int64, image *bytes.Buffer) (usecaseErr error, serverErr error)
	// DeleteImageTodo(todoID int64) (usecaseErr error, serverErr error)

	CreateStatusTodo(name string) (statusTodo *StatusTodo, usecaseErr error, serverErr error)
	GetStatusTodo(id int64) (statusTodo *StatusTodo, usecaseErr error, serverErr error)
}

var (
	ErrTitleIsLong             = errors.New("title is too long")
	ErrDescriptionIsLong       = errors.New("description is too long")
	ErrStatusNotFound          = errors.New("status todo not found")
	ErrNameStatusTodoIsSmall   = errors.New("name status is small")
	ErrStatusTodoAlreadyExists = errors.New("status todo already exists")
	ErrStatusTodoIdNegative    = errors.New("status todo id should to be positive")
)

type DBTodoUsecase struct {
	todoRepository TodoRepository
}

func NewTodoUsecase(todoRepository TodoRepository) TodoUsecase {
	return &DBTodoUsecase{todoRepository}
}

func (usecase *DBTodoUsecase) CreateTodo(title, description string, statusID int64) (todo *Todo, usecaseErr error, serverErr error) {
	if len(title) > 255 {
		usecaseErr = ErrTitleIsLong
		return
	}
	if len(description) > 255 {
		usecaseErr = ErrDescriptionIsLong
		return
	}
	if statusID <= 0 {
		usecaseErr = ErrStatusTodoIdNegative
		return
	}

	statusFound, err := usecase.todoRepository.GetStatusTodo(statusID)
	if err != nil {
		serverErr = err
		return
	}
	if statusFound == nil {
		usecaseErr = ErrStatusNotFound
		return
	}

	todo, err = usecase.todoRepository.InsertTodo(title, description, statusID)
	if err != nil {
		serverErr = err
		return
	}

	return
}

func (usecase *DBTodoUsecase) GetAllTodo() (todos []*Todo, usecaseErr error, serverErr error) {
	todos, serverErr = usecase.todoRepository.GetAllTodo()
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
