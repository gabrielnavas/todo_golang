package controllers

import (
	"api/modules/users/dto"
	"api/modules/users/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController interface {
	Login() func(c *gin.Context)
	Logout() func(c *gin.Context)
}

type LoginControllerGin struct {
	loginUsecase usecases.LoginUsecase
}

func NewLoginController(loginUsecase usecases.LoginUsecase) LoginController {
	return &LoginControllerGin{loginUsecase}
}

func (controller *LoginControllerGin) Login() func(c *gin.Context) {
	return func(c *gin.Context) {
		var body dto.LoginBody
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

		loginResponse, usecaseErr, serverErr := controller.loginUsecase.Login(body.Username, body.Password)
		if serverErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}

		c.JSON(http.StatusCreated, loginResponse)
	}
}

func (controller *LoginControllerGin) Logout() func(c *gin.Context) {
	return func(c *gin.Context) {}
}
