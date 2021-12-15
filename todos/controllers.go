package todos

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TodoController interface {
	CreateTodo() func(c *gin.Context)
	UpdateTodo() func(c *gin.Context)
	DeleteTodo() func(c *gin.Context)
	GetTodo() func(c *gin.Context)
	GetAllTodo() func(c *gin.Context)

	GetImageTodo() func(c *gin.Context)
	UpdateImageTodo() func(c *gin.Context)
	DeleteImageTodo() func(c *gin.Context)

	CreateStatusTodo() func(c *gin.Context)
	UpdateStatusTodo() func(c *gin.Context)
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
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing body"})
			return
		}
		err := body.Validate()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		body.ProcessData()

		todoCreated, usecaseErr, serverErr := controller.todoUsecase.CreateTodo(body.Title, body.Description, body.StatusID)
		if serverErr != nil {
			fmt.Println(serverErr)
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

		var body UpdateTodoBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing body"})
			return
		}
		err = body.Validate()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		body.ProcessData()

		usecaseErr, serverErr := controller.todoUsecase.UpdateTodo(id, body.Title, body.Description, body.StatusID)
		if serverErr != nil {
			fmt.Println(serverErr)
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

func (controller *TodoControllerGin) DeleteTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
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

		usecaseErr, serverErr := controller.todoUsecase.DeleteTodo(id)
		if serverErr != nil {
			fmt.Println(serverErr)
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

func (controller *TodoControllerGin) GetTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
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

		todoFound, usecaseErr, serverErr := controller.todoUsecase.GetTodo(id)
		if serverErr != nil {
			fmt.Println(serverErr)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}
		if todoFound == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "todo not found"})
			return
		}

		c.JSON(http.StatusOK, todoFound)
	}
}

func (controller *TodoControllerGin) GetAllTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		todos, usecaseErr, serverErr := controller.todoUsecase.GetAllTodo()
		if serverErr != nil {
			fmt.Println(serverErr)
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

func (controller *TodoControllerGin) CreateStatusTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		var body CreateStatusTodoBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing body"})
			return
		}
		err := body.Validate()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		statusTodoCreated, usecaseErr, serverErr := controller.todoUsecase.CreateStatusTodo(body.Name)
		if serverErr != nil {
			fmt.Println(serverErr)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}

		c.JSON(http.StatusCreated, statusTodoCreated)
	}
}

func (controller *TodoControllerGin) UpdateStatusTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		idStr, hasId := c.Params.Get("id")
		if !hasId {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing status todo id on url param"})
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing status todo id integer on url param"})
			return
		}

		var body CreateStatusTodoBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing body"})
			return
		}

		body.ProcessData()

		usecaseErr, serverErr := controller.todoUsecase.UpdateStatusTodo(id, body.Name)
		if serverErr != nil {
			fmt.Println(serverErr)
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

func (controller *TodoControllerGin) GetStatusTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		idStr, hasId := c.Params.Get("id")
		if !hasId {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing status todo id on url param"})
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing status todo id integer on url param"})
			return
		}

		statusTodoFound, usecaseErr, serverErr := controller.todoUsecase.GetStatusTodo(id)
		if serverErr != nil {
			fmt.Println(serverErr)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}

		if statusTodoFound == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "status todo not found"})
			return
		}

		c.JSON(http.StatusOK, statusTodoFound)
	}
}
