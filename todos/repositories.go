package todos

import (
	"database/sql"
	"errors"
)

type TodoRepository interface {
	InsertTodo(title, description string, statusID int64) (*Todo, error)
	// UpdateTodo(todoID int64, title, description string) error
	// DeleteTodo(todoID int64) error
	// GetTodo(todoID int64) (*Todo, error)
	GetAllTodo() ([]*Todo, error)

	// GetImageTodo(todoID int64) (*bytes.Buffer, error)
	// UpdateImageTodo(todoID int64, image *bytes.Buffer) error
	// DeleteImageTodo(todoID int64) error

	InsertStatusTodo(name string) (*StatusTodo, error)
	GetStatusTodo(statusID int64) (*StatusTodo, error)
	GetStatusTodoByName(name string) (*StatusTodo, error)
}

type TodoRepositoryPG struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &TodoRepositoryPG{db}
}

func (repo *TodoRepositoryPG) InsertTodo(title, description string, statusID int64) (*Todo, error) {
	var todo Todo
	sqlInsert := `
		INSERT INTO todos.todo (title, description, tstts_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at;
	`
	args := []interface{}{title, description, statusID}
	row := repo.db.QueryRow(sqlInsert, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	err := row.Scan(&todo.ID, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	todo.Title = title
	todo.Description = description
	todo.StatusID = statusID
	return &todo, nil
}

func (repo *TodoRepositoryPG) GetAllTodo() ([]*Todo, error) {
	var todos []*Todo
	sqlGet := `
		SELECT id, title, description, created_at, updated_at, tstts_id
		FROM todos.todo
	`
	rows, err := repo.db.Query(sqlGet)
	if err != nil {
		return nil, nil
	}
	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt, &todo.StatusID)
		if err != nil {
			return nil, nil
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

func (repo *TodoRepositoryPG) InsertStatusTodo(name string) (*StatusTodo, error) {
	var statusTodo StatusTodo
	sqlInsert := `
		INSERT INTO todos.todo_status (name)
		VALUES ($1)
		RETURNING id, created_at, updated_at;
	`
	row := repo.db.QueryRow(sqlInsert, name)
	if row.Err() != nil {
		return nil, row.Err()
	}
	row.Scan(
		&statusTodo.ID,
		&statusTodo.CreatedAt,
		&statusTodo.UpdatedAt,
	)
	statusTodo.Name = name
	return &statusTodo, nil
}

func (repo *TodoRepositoryPG) GetStatusTodo(statusID int64) (*StatusTodo, error) {
	var statusTodo StatusTodo
	sqlGet := `
		SELECT id, name, created_at, updated_at
		from todos.todo_status
		where id=$1;
	`
	row := repo.db.QueryRow(sqlGet, statusID)
	if row.Err() != nil {
		return nil, row.Err()
	}
	err := row.Scan(&statusTodo.ID, &statusTodo.Name, &statusTodo.CreatedAt, &statusTodo.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &statusTodo, nil
}

func (repo *TodoRepositoryPG) GetStatusTodoByName(name string) (*StatusTodo, error) {
	var statusTodo StatusTodo
	sqlGet := `
		SELECT id, name, created_at, updated_at
		FROM todos.todo_status
		WHERE name=$1;
	`

	row := repo.db.QueryRow(sqlGet, name)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&statusTodo.ID,
		&statusTodo.Name,
		&statusTodo.CreatedAt,
		&statusTodo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &statusTodo, nil
}
