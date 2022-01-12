package models

import (
	"bytes"
	"errors"
	"net/mail"
	"time"
)

var (
	ErrNameIsSmall     = errors.New("name is small")
	ErrNameIsLarge     = errors.New("name is large")
	ErrUsernameIsSmall = errors.New("username is small")
	ErrUsernameIsLarge = errors.New("username is large")
	ErrEmailIsInvalid  = errors.New("email is invalid")
	ErrPasswordIsSmall = errors.New("password is small")
	ErrPasswordIsLarge = errors.New("password is large")
)

type User struct {
	ID          int64
	Name        string
	Username    string
	Password    string
	Email       string
	LevelAccess LevelAccess
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Photo       bytes.Buffer
}

type UserSafeHttp struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserRepository interface {
	InsertUser(name, username, password, email string, levelAccess LevelAccess) (*User, error)
	UpdateUser(id int64, name, username, password, email string, levelAccess LevelAccess) error
	DeleteUser(id int64) error
	GetUser(id int64) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetAllUser() ([]*User, error)
	CountUser() (int64, error)

	UpdatePhotoUser(userId int64, photo *bytes.Buffer) error
	GetPhotoUser(id int64) (*bytes.Buffer, error)
}

func NewUser(
	id int64,
	name string,
	username string,
	password string,
	email string,
	levelAccess LevelAccess,
	createdAt time.Time,
	updatedAt time.Time,
	photo bytes.Buffer,
) (*User, error) {
	user := &User{id, name, username, password, email, levelAccess, createdAt, updatedAt, photo}
	err := user.Valid()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Valid() error {
	if len(u.Name) < 2 {
		return ErrNameIsSmall
	}
	if len(u.Name) > 255 {
		return ErrNameIsLarge
	}
	if len(u.Username) < 2 {
		return ErrUsernameIsSmall
	}
	if len(u.Username) > 255 {
		return ErrUsernameIsLarge
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return ErrEmailIsInvalid
	}
	if err := u.LevelAccess.Valid(); err != nil {
		return err
	}
	if len(u.Password) < 6 {
		return ErrPasswordIsSmall
	}
	if len(u.Password) > 255 {
		return ErrPasswordIsLarge
	}
	return nil
}

func (u *User) ToSafeHttp() *UserSafeHttp {
	return &UserSafeHttp{
		ID:        u.ID,
		Name:      u.Name,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
