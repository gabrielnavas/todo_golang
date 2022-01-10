package users

import (
	"bytes"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type UserRepositoryPG struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryPG{db}
}

func (repo *UserRepositoryPG) InsertUser(name, username, password, email string) (*User, error) {
	var user User
	sqlInsert := `
		INSERT INTO users.user (name, email, username, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at;
	`
	args := []interface{}{name, email, username, password}
	row := repo.db.QueryRow(sqlInsert, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	user.Name = name
	user.Email = email
	user.Username = username
	return &user, nil
}

func (repo *UserRepositoryPG) UpdateUser(id int64, name, username, password, email string) error {
	now := time.Now().UTC()
	sqlUpdate := `
		UPDATE users.user
		SET 
			name=$2,
			email=$3,
			username=$4,
			password=$5,
			updated_at=$6
		WHERE id=$1
	`
	args := []interface{}{id, name, email, username, password, now}
	_, err := repo.db.Exec(sqlUpdate, args...)
	return err
}

func (repo *UserRepositoryPG) DeleteUser(id int64) error {
	sqlDelete := `
		DELETE FROM users.user
		WHERE id=$1;
	`
	_, err := repo.db.Exec(sqlDelete, id)
	return err
}

func (repo *UserRepositoryPG) GetUser(id int64) (*User, error) {
	var user User
	var bufferImage = []byte{}

	sqlGet := `
		SELECT id, name, email, username, password, created_at, updated_at, photo
		FROM users.user
		WHERE id=$1;
	`

	row := repo.db.QueryRow(sqlGet, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
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
		user.Photo.ReadFrom(reader)
	}

	return &user, nil
}

func (repo *UserRepositoryPG) GetUserByEmail(email string) (*User, error) {
	var user User
	var bufferImage = []byte{}

	sqlGet := `
		SELECT id, name, email, username, password, created_at, updated_at, photo
		FROM users.user
		WHERE LOWER(email)=$1;
	`

	emailLower := strings.ToLower(email)
	row := repo.db.QueryRow(sqlGet, emailLower)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
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
		user.Photo.ReadFrom(reader)
	}

	return &user, nil
}

func (repo *UserRepositoryPG) GetUserByUsername(username string) (*User, error) {
	var user User
	var bufferImage = []byte{}

	sqlGet := `
		SELECT id, name, email, username, password, created_at, updated_at, photo
		FROM users.user
		WHERE LOWER(username)=$1;
	`

	emailLower := strings.ToLower(username)
	row := repo.db.QueryRow(sqlGet, emailLower)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
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
		user.Photo.ReadFrom(reader)
	}

	return &user, nil
}

func (repo *UserRepositoryPG) GetAllUser() ([]*User, error) {
	var todos = make([]*User, 0)
	sqlGet := `
		SELECT id, name, email, username, password, created_at, updated_at, photo
		FROM users.user
		ORDER BY created_at DESC;
	`
	rows, err := repo.db.Query(sqlGet)
	if err != nil {
		return nil, nil
	}

	for rows.Next() {
		var user User
		var bufferImage = []byte{}

		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
			&bufferImage,
		)
		if err != nil {
			return nil, nil
		}

		if len(bufferImage) > 0 {
			reader := bytes.NewReader(bufferImage)
			user.Photo.ReadFrom(reader)
		}

		todos = append(todos, &user)
	}

	return todos, nil
}

func (repo *UserRepositoryPG) UpdatePhotoUser(userId int64, photo *bytes.Buffer) error {
	var imageToArgs interface{} = nil
	if photo != nil {
		imageToArgs = photo.Bytes()
	}
	sqlUpdate := `
		UPDATE users.user
		SET photo=$2
		WHERE id=$1;
	`
	args := []interface{}{userId, imageToArgs}
	_, err := repo.db.Exec(sqlUpdate, args...)
	return err
}

func (repo *UserRepositoryPG) GetPhotoUser(id int64) (*bytes.Buffer, error) {
	buffImage := []byte{}
	sqlGet := `
		SELECT photo 
		FROM users.user
		WHERE id=$1;
	`
	row := repo.db.QueryRow(sqlGet, id)
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
