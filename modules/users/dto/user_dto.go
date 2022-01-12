package dto

import (
	"api/modules/users/models"
	"bytes"
	"errors"
	"mime/multipart"
	"strings"
)

type ChangePasswordBody struct {
	OldPassword             string `json:"oldPassword"`
	NewPassword             string `json:"newPassword"`
	NewPasswordConfirmation string `json:"newPasswordConfirmation"`
}

func (body *ChangePasswordBody) Validate() error {
	if body.OldPassword == "" {
		return errors.New("old password is empty")
	}
	if body.NewPassword == "" {
		return errors.New("new password is empty")
	}
	if body.NewPasswordConfirmation == "" {
		return errors.New("new password confirmation is empty")
	}
	return nil
}

func (body *ChangePasswordBody) ProcessData() {
	body.OldPassword = strings.TrimSpace(body.OldPassword)
	body.NewPassword = strings.TrimSpace(body.NewPassword)
	body.NewPasswordConfirmation = strings.TrimSpace(body.NewPasswordConfirmation)
}

type CreateUserBody struct {
	Name                 string `json:"name"`
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
	Email                string `json:"email"`
}

func (body *CreateUserBody) Validate() error {
	if body.Name == "" {
		return errors.New("name is empty")
	}
	if body.Username == "" {
		return errors.New("username is empty")
	}
	if body.Password == "" {
		return errors.New("password is empty")
	}
	if body.PasswordConfirmation == "" {
		return errors.New("password confirmation is empty")
	}
	if body.Email == "" {
		return errors.New("email is empty")
	}
	return nil
}

func (body *CreateUserBody) ProcessData() {
	body.Name = strings.TrimSpace(body.Name)
	body.Username = strings.TrimSpace(body.Username)
	body.Password = strings.TrimSpace(body.Password)
	body.Email = strings.TrimSpace(body.Email)
}

type UpdatePhotoUserDTO struct {
	UserId      int64
	Size        int64
	ContentType string
	BufferFile  bytes.Buffer
}

func NewUpdatePhotoUserDTO(userId int64, fileHeader *multipart.FileHeader) (*UpdatePhotoUserDTO, error) {
	// open file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create instance struct to return
	dto := UpdatePhotoUserDTO{}
	dto.UserId = userId

	// get buffer
	sizeBytes, err := dto.BufferFile.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	//set size bytes
	dto.Size = sizeBytes

	// handle content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		return nil, errors.New("Content-Type header is empty")
	}
	dto.ContentType = contentType

	return &dto, nil
}

func (file *UpdatePhotoUserDTO) Validate() error {
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

func (file *UpdatePhotoUserDTO) checkContentType() error {
	contentTypesExpected := []string{"image/jpeg", "image/png"}
	for _, contentType := range contentTypesExpected {
		if contentType == file.ContentType {
			return nil
		}
	}
	return errors.New("unsupported content type")
}

func (file *UpdatePhotoUserDTO) checkSize() error {
	fiveMB := int64(1024 * 1024 * 5)
	if file.Size > fiveMB {
		return errors.New("image is more that 5MB")
	}
	return nil
}

type UpdateUserBody struct {
	CreateUserBody
	LevelAccess models.LevelAccess `json:"levelAccess"`
}
