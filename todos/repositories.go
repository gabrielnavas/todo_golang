package todos

import (
	"bytes"
	"database/sql"
	"errors"
	"time"
)

type TodoRepository interface {
	InsertTodo(title, description string, statusID int64) (*Todo, error)
	UpdateTodo(todoID int64, title, description string, statusTodoId int64) error
	DeleteTodo(todoID int64) error
	GetTodo(todoID int64) (*Todo, error)
	GetAllTodo() ([]*Todo, error)
	CountTodoByStatus(statusTodoId int64) (int64, error)

	UpdateImageTodo(todoID int64, image *bytes.Buffer) error
	GetImageTodo(todoID int64) (*bytes.Buffer, error)

	InsertStatusTodo(name string) (*StatusTodo, error)
	UpdateStatusTodo(id int64, name string) error
	GetAllStatusTodo() ([]*StatusTodo, error)
	GetStatusTodo(statusID int64) (*StatusTodo, error)
	GetStatusTodoByName(name string) (*StatusTodo, error)
	DeleteStatusTodo(id int64) error
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

func (repo *TodoRepositoryPG) UpdateTodo(todoID int64, title, description string, statusTodoId int64) error {
	now := time.Now().UTC()
	sqlUpdate := `
		UPDATE todos.todo
		SET 
			title=$2,
			description=$3,
			tstts_id=$4,
			updated_at=$5
		WHERE id=$1
	`
	args := []interface{}{todoID, title, description, statusTodoId, now}
	_, err := repo.db.Exec(sqlUpdate, args...)
	return err
}

func (repo *TodoRepositoryPG) DeleteTodo(todoID int64) error {
	sqlDelete := `
		DELETE FROM todos.todo
		WHERE id=$1;
	`
	_, err := repo.db.Exec(sqlDelete, todoID)
	return err
}

func (repo *TodoRepositoryPG) GetTodo(todoID int64) (*Todo, error) {
	var todo Todo
	var bufferImage = []byte{}

	sqlGet := `
		SELECT id, title, description, created_at, updated_at, tstts_id, image
		FROM todos.todo
		WHERE id=$1;
	`

	row := repo.db.QueryRow(sqlGet, todoID)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.StatusID,
		&bufferImage,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if len(bufferImage) > 0 {
		reader := bytes.NewReader(bufferImage)
		todo.Image.ReadFrom(reader)
	}

	return &todo, nil
}

func (repo *TodoRepositoryPG) GetAllTodo() ([]*Todo, error) {
	var todos = make([]*Todo, 0)
	sqlGet := `
		SELECT id, title, description, created_at, updated_at, tstts_id, image
		FROM todos.todo
	`
	rows, err := repo.db.Query(sqlGet)
	if err != nil {
		return nil, nil
	}

	for rows.Next() {
		var todo Todo
		var bufferImage = []byte{}

		err = rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.StatusID,
			&bufferImage,
		)
		if err != nil {
			return nil, nil
		}

		if len(bufferImage) > 0 {
			reader := bytes.NewReader(bufferImage)
			todo.Image.ReadFrom(reader)
		}

		todos = append(todos, &todo)
	}

	return todos, nil
}

func (repo *TodoRepositoryPG) CountTodoByStatus(statusTodoId int64) (int64, error) {
	var count int64
	sqlCount := `
		SELECT COUNT(*)
		FROM todos.todo
		WHERE tstts_id=$1
	`
	row := repo.db.QueryRow(sqlCount, statusTodoId)
	if row.Err() != nil {
		return -1, row.Err()
	}
	err := row.Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (repo *TodoRepositoryPG) UpdateImageTodo(todoID int64, image *bytes.Buffer) error {
	var imageToArgs interface{}
	if image == nil {
		imageToArgs = nil
	} else {
		imageToArgs = image.Bytes()
	}
	sqlUpdate := `
		UPDATE todos.todo
		SET image=$2
		WHERE id=$1;
	`
	args := []interface{}{todoID, imageToArgs}
	_, err := repo.db.Exec(sqlUpdate, args...)
	return err
}

func (repo *TodoRepositoryPG) GetImageTodo(todoID int64) (*bytes.Buffer, error) {
	buffImage := []byte{}
	sqlGet := `
		SELECT image 
		FROM todos.todo
		WHERE id=$1;
	`
	row := repo.db.QueryRow(sqlGet, todoID)
	if row.Err() != nil {
		return nil, row.Err()
	}
	err := row.Scan(&buffImage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return bytes.NewBuffer(buffImage), nil
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

func (repo *TodoRepositoryPG) UpdateStatusTodo(id int64, name string) error {
	sqlUpdate := `
		UPDATE todos.todo_status
		SET name=$2
		WHERE id=$1;; 
	`
	args := []interface{}{id, name}
	_, error := repo.db.Exec(sqlUpdate, args...)
	return error
}

func (repo *TodoRepositoryPG) GetAllStatusTodo() ([]*StatusTodo, error) {
	var allStatusTodo = make([]*StatusTodo, 0)
	sqlGet := `
		SELECT id, name, created_at, updated_at 
		FROM todos.todo_status;
	`
	rows, err := repo.db.Query(sqlGet)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var statusTodo StatusTodo
		err := rows.Scan(&statusTodo.ID, &statusTodo.Name, &statusTodo.CreatedAt, &statusTodo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		allStatusTodo = append(allStatusTodo, &statusTodo)
	}
	return allStatusTodo, nil
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

func (repo *TodoRepositoryPG) DeleteStatusTodo(id int64) error {
	sqlDelete := `
		DELETE FROM todos.todo_status
		WHERE id=$1;
	`
	_, error := repo.db.Exec(sqlDelete, id)
	return error
}
