package todos

import (
	"api/modules/users/middlewares"
	"api/modules/users/models"
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
	GetAllTodos() func(c *gin.Context)

	GetImageTodo() func(c *gin.Context)
	UpdateImageTodo() func(c *gin.Context)
	DeleteImageTodo() func(c *gin.Context)

	CreateStatusTodo() func(c *gin.Context)
	UpdateStatusTodo() func(c *gin.Context)
	GetStatusTodo() func(c *gin.Context)
	GetAllStatusTodo() func(c *gin.Context)
	DeleteStatusTodo() func(c *gin.Context)
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

		// get user id
		userIdValue, exists := c.Get(middlewares.UserId)
		if !exists {
			fmt.Println("need user id for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		userId := userIdValue.(int64)

		// get level access
		levelAccessValue, exists := c.Get(middlewares.LevelAccess)
		if !exists {
			fmt.Println("need level access for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		levelAccess := levelAccessValue.(models.LevelAccess)

		// check authorization
		if levelAccess < models.BasicLevelAccess {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have authorization"})
			return
		}

		todoCreated, usecaseErr, serverErr := controller.todoUsecase.CreateTodo(body.Title, body.Description, body.StatusID, userId)
		if serverErr != nil {
			fmt.Println(serverErr)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}
		c.JSON(http.StatusCreated, todoCreated.ToDtoHttpResponse())
	}
}

func (controller *TodoControllerGin) GetImageTodo() func(c *gin.Context) {
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

		// get level access
		levelAccessValue, exists := c.Get(middlewares.LevelAccess)
		if !exists {
			fmt.Println("need level access for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		levelAccess := levelAccessValue.(models.LevelAccess)

		// check authorization
		if levelAccess < models.BasicLevelAccess {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have authorization"})
			return
		}

		bufferImage, usecaseErr, serverErr := controller.todoUsecase.GetImageTodo(id)
		if serverErr != nil {
			fmt.Println(serverErr)
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
		c.Data(http.StatusOK, contentType, bufferImage.Bytes())
	}
}

func (controller *TodoControllerGin) UpdateImageTodo() func(c *gin.Context) {
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

		// get level access
		levelAccessValue, exists := c.Get(middlewares.LevelAccess)
		if !exists {
			fmt.Println("need level access for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		levelAccess := levelAccessValue.(models.LevelAccess)

		// check authorization
		if levelAccess < models.BasicLevelAccess {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have authorization"})
			return
		}

		// get image
		fileHeader, _ := c.FormFile("image")
		updateImageTodoFile, err := NewUpdateImageTodoFile(id, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		err = updateImageTodoFile.Validate()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// update
		usecaseErr, serverErr := controller.todoUsecase.UpdateImageTodo(updateImageTodoFile)
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

func (controller *TodoControllerGin) DeleteImageTodo() func(c *gin.Context) {
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

		// get level access
		levelAccessValue, exists := c.Get(middlewares.LevelAccess)
		if !exists {
			fmt.Println("need level access for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		levelAccess := levelAccessValue.(models.LevelAccess)

		// check authorization
		if levelAccess < models.BasicLevelAccess {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have authorization"})
			return
		}

		usecaseErr, serverErr := controller.todoUsecase.DeleteImageTodo(id)
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

		// get user id
		userIdValue, exists := c.Get(middlewares.UserId)
		if !exists {
			fmt.Println("need user id for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		userId := userIdValue.(int64)

		// get level access
		levelAccessValue, exists := c.Get(middlewares.LevelAccess)
		if !exists {
			fmt.Println("need level access for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		levelAccess := levelAccessValue.(models.LevelAccess)

		// check authorization
		if levelAccess < models.BasicLevelAccess {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have authorization"})
			return
		}

		usecaseErr, serverErr := controller.todoUsecase.UpdateTodo(id, body.Title, body.Description, body.StatusID, userId)
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

		// get level access
		levelAccessValue, exists := c.Get(middlewares.LevelAccess)
		if !exists {
			fmt.Println("need level access for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		levelAccess := levelAccessValue.(models.LevelAccess)

		// check authorization
		if levelAccess < models.BasicLevelAccess {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have authorization"})
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

		// get level access
		levelAccessValue, exists := c.Get(middlewares.LevelAccess)
		if !exists {
			fmt.Println("need level access for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		levelAccess := levelAccessValue.(models.LevelAccess)

		// check authorization
		if levelAccess < models.BasicLevelAccess {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have authorization"})
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

		c.JSON(http.StatusOK, todoFound.ToDtoHttpResponse())
	}
}

func (controller *TodoControllerGin) GetAllTodos() func(c *gin.Context) {
	return func(c *gin.Context) {
		// get level access
		levelAccessValue, exists := c.Get(middlewares.LevelAccess)
		if !exists {
			fmt.Println("need level access for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		levelAccess := levelAccessValue.(models.LevelAccess)

		// check authorization
		if levelAccess < models.BasicLevelAccess {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have authorization"})
			return
		}

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

		var wrappedResponse []*TodoDtoHttpResponse = make([]*TodoDtoHttpResponse, 0)
		for _, todo := range todos {
			wrappedResponse = append(wrappedResponse, todo.ToDtoHttpResponse())
		}

		c.JSON(http.StatusOK, wrappedResponse)
	}
}

func (controller *TodoControllerGin) CreateStatusTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		// get body
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

		// get user id
		userIdValue, exists := c.Get(middlewares.UserId)
		if !exists {
			fmt.Println("need user id for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		userId := userIdValue.(int64)

		// get level access
		levelAccessValue, exists := c.Get(middlewares.LevelAccess)
		if !exists {
			fmt.Println("need level access for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		levelAccess := levelAccessValue.(models.LevelAccess)

		// check authorization
		if levelAccess < models.BasicLevelAccess {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have authorization"})
			return
		}

		// create status todo
		statusTodoCreated, usecaseErr, serverErr := controller.todoUsecase.CreateStatusTodo(body.Name, userId)
		if serverErr != nil {
			fmt.Println(serverErr)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}

		//return
		c.JSON(http.StatusCreated, statusTodoCreated)
	}
}

// TODO: refatorar para vir com userId
func (controller *TodoControllerGin) UpdateStatusTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		idStr, hasId := c.Params.Get("id")
		if !hasId {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing status todo id on url param"})
			return
		}
		statusId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing status todo id integer on url param"})
			return
		}

		// TODO: mudar para updatestatusTodobody
		var body CreateStatusTodoBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing body"})
			return
		}

		body.ProcessData()

		// get user id
		userIdValue, exists := c.Get(middlewares.UserId)
		if !exists {
			fmt.Println("need user id for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		userId := userIdValue.(int64)

		// get level access
		levelAccessValue, exists := c.Get(middlewares.LevelAccess)
		if !exists {
			fmt.Println("need level access for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		levelAccess := levelAccessValue.(models.LevelAccess)

		// check authorization
		if levelAccess < models.BasicLevelAccess {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have authorization"})
			return
		}

		usecaseErr, serverErr := controller.todoUsecase.UpdateStatusTodo(userId, statusId, body.Name)
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

func (controller *TodoControllerGin) GetAllStatusTodo() func(c *gin.Context) {
	return func(c *gin.Context) {
		// get level access
		levelAccessValue, exists := c.Get(middlewares.LevelAccess)
		if !exists {
			fmt.Println("need level access for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		levelAccess := levelAccessValue.(models.LevelAccess)

		// check authorization
		if levelAccess < models.BasicLevelAccess {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have authorization"})
			return
		}

		allStatusTodo, usecaseErr, serverErr := controller.todoUsecase.GetAllStatusTodo()
		if serverErr != nil {
			fmt.Println(serverErr)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		if usecaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": usecaseErr.Error()})
			return
		}
		c.JSON(http.StatusOK, allStatusTodo)
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

		// get user id
		userIdValue, exists := c.Get(middlewares.UserId)
		if !exists {
			fmt.Println("need user id for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		userId := userIdValue.(int64)

		// get level access
		levelAccessValue, exists := c.Get(middlewares.LevelAccess)
		if !exists {
			fmt.Println("need level access for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		levelAccess := levelAccessValue.(models.LevelAccess)

		// check authorization
		if levelAccess < models.BasicLevelAccess {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have authorization"})
			return
		}

		statusTodoFound, usecaseErr, serverErr := controller.todoUsecase.GetStatusTodo(userId, id)
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

func (controller *TodoControllerGin) DeleteStatusTodo() func(c *gin.Context) {
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

		// get user id
		userIdValue, exists := c.Get(middlewares.UserId)
		if !exists {
			fmt.Println("need user id for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		userId := userIdValue.(int64)

		// get level access
		levelAccessValue, exists := c.Get(middlewares.LevelAccess)
		if !exists {
			fmt.Println("need level access for create status todo")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		levelAccess := levelAccessValue.(models.LevelAccess)

		// check authorization
		if levelAccess < models.BasicLevelAccess {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "you don't have authorization"})
			return
		}

		usecaseErr, serverErr := controller.todoUsecase.DeleteStatusTodo(userId, id)
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
