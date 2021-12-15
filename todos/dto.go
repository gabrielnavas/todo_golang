package todos

import (
	"errors"
	"strings"
)

type CreateTodoBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusID    int64  `json:"statusId"`
}

func (body *CreateTodoBody) Validate() error {
	if body.Title == "" {
		return errors.New("missing title")
	}
	if body.Description == "" {
		return errors.New("missing description")
	}
	return nil
}

func (body *CreateTodoBody) ProcessData() {
	body.Title = strings.TrimSpace(body.Title)
	body.Description = strings.TrimSpace(body.Description)
}

type UpdateTodoBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusID    int64  `json:"statusId"`
}

func (body *UpdateTodoBody) Validate() error {
	if body.Title == "" {
		return errors.New("missing title")
	}
	if body.Description == "" {
		return errors.New("missing description")
	}
	return nil
}

func (body *UpdateTodoBody) ProcessData() {
	body.Title = strings.TrimSpace(body.Title)
	body.Description = strings.TrimSpace(body.Description)
}

type CreateStatusTodoBody struct {
	Name string `json:"name"`
}

func (body *CreateStatusTodoBody) Validate() error {
	if body.Name == "" {
		return errors.New("missing name")
	}
	return nil
}

func (body *CreateStatusTodoBody) ProcessData() {
	body.Name = strings.TrimSpace(body.Name)
}
