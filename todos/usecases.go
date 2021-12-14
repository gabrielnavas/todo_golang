package todos

import (
	"errors"
)

type TodoUsecase interface {
	CreateTodo(title, description string, statusID int64) (todo *Todo, usecaseErr error, serverErr error)
	// GetImageTodo(todoID int64) (image *bytes.Buffer, serverErr error)
	// UpdateImageTodo(todoID int64, image *bytes.Buffer) (usecaseErr error, serverErr error)
	// DeleteImageTodo(todoID int64) (usecaseErr error, serverErr error)
	// UpdateTodo(todoID int64, title, description string) (usecaseErr error, serverErr error)
	// DeleteTodo(todoID int64) (usecaseErr error, serverErr error)
	// GetTodo(todoID int64) (todo *Todo, usecaseErr error, serverErr error)
	GetAllTodo() (todos []*Todo, usecaseErr error, serverErr error)
}

var (
	ErrTitleIsLong       = errors.New("title is too long")
	ErrDescriptionIsLong = errors.New("description is too long")
	ErrStatusNotFound    = errors.New("status todo not found")
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
