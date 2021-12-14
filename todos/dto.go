package todos

import "errors"

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
	if body.StatusID <= 0 {
		return errors.New("missing statusId")
	}
	return nil
}
