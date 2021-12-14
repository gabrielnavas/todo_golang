package todos

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TodoController interface {
	CreateTodo() func(c *gin.Context)
	GetImageTodo() func(c *gin.Context)
	UpdateImageTodo() func(c *gin.Context)
	DeleteImageTodo() func(c *gin.Context)
	UpdateTodo() func(c *gin.Context)
	DeleteTodo() func(c *gin.Context)
	GetTodo() func(c *gin.Context)
	GetAllTodo() func(c *gin.Context)
	GetStatusTodo() func(c *gin.Context)
}

type TodoControllerGin struct {
	todoUsecase TodoUsecase
}

func NewTodoController(todoUsecase TodoUsecase) TodoController {
	return &TodoControllerGin{todoUsecase}
}

func (controller *TodoControllerGin) CreateTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		var body CreateTodoBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		err := body.Validate()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		todoCreated, usecaseErr, serverErr := controller.todoUsecase.CreateTodo(body.Title, body.Description, body.StatusID)
		if serverErr != nil {
			fmt.Println("deu merda", serverErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}

		c.JSON(http.StatusCreated, todoCreated)
	}
}

func (controller *TodoControllerGin) GetImageTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "server error"})
	}
}

func (controller *TodoControllerGin) UpdateImageTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "server error"})
	}
}

func (controller *TodoControllerGin) DeleteImageTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "server error"})
	}
}

func (controller *TodoControllerGin) UpdateTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "server error"})
	}
}

func (controller *TodoControllerGin) DeleteTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "server error"})
	}
}

func (controller *TodoControllerGin) GetTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "server error"})
	}
}

func (controller *TodoControllerGin) GetAllTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		todos, usecaseErr, serverErr := controller.todoUsecase.GetAllTodo()
		if serverErr != nil {
			fmt.Println("deu merda", serverErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}
		c.JSON(http.StatusOK, todos)
	}
}

func (controller *TodoControllerGin) GetStatusTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "server error"})
	}
}
