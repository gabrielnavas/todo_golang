package repositories

import (
	"bytes"
	"database/sql"
	"errors"
	"strings"
	"time"

	"api/modules/users/models"
)

type UserRepositoryPG struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) models.UserRepository {
	return &UserRepositoryPG{db}
}

func (repo *UserRepositoryPG) InsertUser(name, username, password, email string, levelAccess models.LevelAccess) (*models.User, error) {
	var user models.User
	sqlInsert := `
		INSERT INTO users.user (name, email, username, password, level_access)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at;
	`
	args := []interface{}{name, email, username, password, levelAccess}
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
	user.LevelAccess = levelAccess
	return &user, nil
}

func (repo *UserRepositoryPG) UpdateUser(id int64, name, username, password, email string, levelAccess models.LevelAccess) error {
	now := time.Now().UTC()
	sqlUpdate := `
		UPDATE users.user
		SET 
			name=$2,
			email=$3,
			username=$4,
			level_access=$5,
			password=$6,
			updated_at=$7
		WHERE id=$1
	`
	args := []interface{}{id, name, email, username, levelAccess, password, now}
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

func (repo *UserRepositoryPG) GetUser(id int64) (*models.User, error) {
	var user models.User
	var bufferImage = []byte{}

	sqlGet := `
		SELECT id, name, email, username, password, level_access, created_at, updated_at, photo
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
		&user.LevelAccess,
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

func (repo *UserRepositoryPG) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	var bufferImage = []byte{}

	sqlGet := `
		SELECT id, name, email, username, password, level_access, created_at, updated_at, photo
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
		&user.LevelAccess,
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

func (repo *UserRepositoryPG) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	var bufferImage = []byte{}

	sqlGet := `
		SELECT id, name, email, username, password, level_access, created_at, updated_at, photo
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
		&user.LevelAccess,
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

func (repo *UserRepositoryPG) GetAllUser() ([]*models.User, error) {
	var todos = make([]*models.User, 0)
	sqlGet := `
		SELECT id, name, email, username, password, level_access, created_at, updated_at, photo
		FROM users.user
		ORDER BY created_at DESC;
	`
	rows, err := repo.db.Query(sqlGet)
	if err != nil {
		return nil, nil
	}

	for rows.Next() {
		var user models.User
		var bufferImage = []byte{}

		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.LevelAccess,
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

func (repo *UserRepositoryPG) CountUser() (int64, error) {
	var count int64
	sqlCount := `
		SELECT COUNT(*)
		FROM users.user;
	`

	result, err := repo.db.Query(sqlCount)
	if err != nil {
		return -1, err
	}
	result.Next()
	result.Scan(&count)
	return count, nil
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
