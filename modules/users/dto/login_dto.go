package dto

import (
	"api/modules/users/models"
	"errors"
	"strings"
)

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (body *LoginBody) Validate() error {
	if body.Username == "" {
		return errors.New("username is empty")
	}
	if body.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}

func (body *LoginBody) ProcessData() {
	body.Password = strings.TrimSpace(body.Password)
	body.Username = strings.TrimSpace(body.Username)
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type LoginResponseSafe struct {
	Token string              `json:"token"`
	User  models.UserSafeHttp `json:"user"`
}

func (dto *LoginResponse) ToSafeHttp() *LoginResponseSafe {
	return &LoginResponseSafe{
		dto.Token,
		*dto.User.ToSafeHttp(),
	}
}
