package todos

import (
	"bytes"
	"errors"
	"mime/multipart"
	"strings"
)

type CreateTodoBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusID    int64  `json:"statusId"`
	UserId      int64  `json:"userId"`
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
	UserId      int64  `json:"userId"`
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

type UpdateImageTodoDTO struct {
	TodoId      int64
	Size        int64
	ContentType string
	BufferFile  bytes.Buffer
}

func NewUpdateImageTodoFile(todoId int64, fileHeader *multipart.FileHeader) (*UpdateImageTodoDTO, error) {
	// open file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create instance struct to return
	updateImageTodoFile := UpdateImageTodoDTO{}
	updateImageTodoFile.TodoId = todoId

	// get buffer
	sizeBytes, err := updateImageTodoFile.BufferFile.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	//set size bytes
	updateImageTodoFile.Size = sizeBytes

	// handle content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		return nil, errors.New("Content-Type header is empty")
	}
	updateImageTodoFile.ContentType = contentType

	return &updateImageTodoFile, nil
}

func (file *UpdateImageTodoDTO) Validate() error {
	err := file.checkContentType()
	if err != nil {
		return err
	}

	err = file.checkSize()
	if err != nil {
		return err
	}

	return nil
}

func (file *UpdateImageTodoDTO) checkContentType() error {
	contentTypesExpected := []string{"image/jpeg", "image/png"}
	for _, contentType := range contentTypesExpected {
		if contentType == file.ContentType {
			return nil
		}
	}
	return errors.New("unsupported content type")
}

func (file *UpdateImageTodoDTO) checkSize() error {
	fiveMB := int64(1024 * 1024 * 5)
	if file.Size > fiveMB {
		return errors.New("image is more that 5MB")
	}
	return nil
}
