package controllers

import (
	"api/modules/users/dto"
	"api/modules/users/models"
	"api/modules/users/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	CreateUser() func(c *gin.Context)
	UpdateUser() func(c *gin.Context)
	DeleteUser() func(c *gin.Context)
	GetUser() func(c *gin.Context)
	GetAllUser() func(c *gin.Context)
	ChangePassword() func(c *gin.Context)

	UpdatePhotoUser() func(c *gin.Context)
	DeletePhotoUser() func(c *gin.Context)
	GetPhotoUser() func(c *gin.Context)
}

type UserControllerGin struct {
	userUsecase usecases.UserUsecase
}

func NewUserController(userUsecase usecases.UserUsecase) UserController {
	return &UserControllerGin{userUsecase}
}

func (controller *UserControllerGin) CreateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var body dto.CreateUserBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing body"})
			return
		}
		body.ProcessData()
		err := body.Validate()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		userCreated, usecaseErr, serverErr := controller.userUsecase.CreateUser(body.Name, body.Username, body.Password, body.PasswordConfirmation, body.Email)
		if serverErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}

		userSafe := userCreated.ToSafeHttp()
		c.JSON(http.StatusCreated, userSafe)
	}
}

func (controller *UserControllerGin) UpdateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		idStr, hasId := c.Params.Get("id")
		if !hasId {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing user id on url param"})
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing user id integer on url param"})
			return
		}

		var body dto.UpdateUserBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing body"})
			return
		}
		body.ProcessData()
		err = body.Validate()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		usecaseErr, serverErr := controller.userUsecase.UpdateUser(id, body.Name, body.Username, body.Password, body.Email, body.LevelAccess)
		if serverErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func (controller *UserControllerGin) DeleteUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		idStr, hasId := c.Params.Get("id")
		if !hasId {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing user id on url param"})
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing user id integer on url param"})
			return
		}

		usecaseErr, serverErr := controller.userUsecase.DeleteUser(id)
		if serverErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func (controller *UserControllerGin) GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		idStr, hasId := c.Params.Get("id")
		if !hasId {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing user id on url param"})
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing user id integer on url param"})
			return
		}

		userFound, usecaseErr, serverErr := controller.userUsecase.GetUser(id)
		if serverErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}
		if userFound == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "todo not found"})
			return
		}

		userSafe := userFound.ToSafeHttp()
		c.JSON(http.StatusOK, userSafe)
	}
}

func (controller *UserControllerGin) GetAllUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		users, usecaseErr, serverErr := controller.userUsecase.GetAllUser()
		if serverErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}

		var wrappedUsersResponse []*models.UserSafeHttp = make([]*models.UserSafeHttp, 0)
		for _, user := range users {
			wrappedUsersResponse = append(wrappedUsersResponse, user.ToSafeHttp())
		}

		c.JSON(http.StatusOK, wrappedUsersResponse)
	}
}

func (controller *UserControllerGin) ChangePassword() func(c *gin.Context) {
	return func(c *gin.Context) {
		idStr, hasId := c.Params.Get("id")
		if !hasId {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing user id on url param"})
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing user id integer on url param"})
			return
		}

		var body dto.ChangePasswordBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing body"})
			return
		}
		body.ProcessData()
		err = body.Validate()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		usecaseErr, serverErr := controller.userUsecase.ChangePassword(id, body.OldPassword, body.NewPassword, body.NewPasswordConfirmation)
		if serverErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func (controller *UserControllerGin) UpdatePhotoUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		// get id
		idStr, hasId := c.Params.Get("id")
		if !hasId {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing id id on url param"})
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing id id integer on url param"})
			return
		}

		// get photo
		fileHeader, _ := c.FormFile("photo")
		userPhotoDto, err := dto.NewUpdatePhotoUserDTO(id, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		err = userPhotoDto.Validate()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// update
		usecaseErr, serverErr := controller.userUsecase.UpdatePhotoUser(userPhotoDto)
		if serverErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func (controller *UserControllerGin) DeletePhotoUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		idStr, hasId := c.Params.Get("id")
		if !hasId {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing user id on url param"})
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing user id integer on url param"})
			return
		}

		usecaseErr, serverErr := controller.userUsecase.DeletePhotoUser(id)
		if serverErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func (controller *UserControllerGin) GetPhotoUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		// get id
		idStr, hasId := c.Params.Get("id")
		if !hasId {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing todo id on url param"})
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing todo id integer on url param"})
			return
		}

		bufferImage, usecaseErr, serverErr := controller.userUsecase.GetPhotoUser(id)
		if serverErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}
		if bufferImage.Len() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "image not found"})
			return
		}

		contentType := http.DetectContentType(bufferImage.Bytes())
		bytesFromImage := bufferImage.Bytes()
		c.Data(http.StatusOK, contentType, bytesFromImage)
	}
}
