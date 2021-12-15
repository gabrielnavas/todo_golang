package main

import (
	"api/database"
	"api/env"
	"api/todos"

	"github.com/gin-gonic/gin"
)

func main() {
	// env
	env, err := env.NewEnvironment()
	if err != nil {
		panic(err)
	}

	// database
	db, err := database.MakeConnection(
		env.Database.User,
		env.Database.Host,
		env.Database.Port,
		env.Database.Password,
		env.Database.Dbname,
		env.Database.Sslmode,
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// routers
	router := gin.Default()
	{
		// todos
		repo := todos.NewTodoRepository(db)
		todoUsecase := todos.NewTodoUsecase(repo)
		controller := todos.NewTodoController(todoUsecase)
		router.POST("/todos", controller.CreateTodo())
		router.GET("/todos/:id", controller.GetTodo())
		router.GET("/todos", controller.GetAllTodos())
		router.PUT("/todos/:id", controller.UpdateTodo())
		router.DELETE("/todos/:id", controller.DeleteTodo())

		// todo status
		router.POST("todos/status", controller.CreateStatusTodo())
		router.GET("todos/status/:id", controller.GetStatusTodo())
		router.GET("todos/status", controller.GetAllStatusTodo())
		router.PUT("todos/status/:id", controller.UpdateStatusTodo())
		router.DELETE("todos/status/:id", controller.DeleteStatusTodo())
	}

	router.Run()

}
