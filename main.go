package main

import (
	"api/database"
	"api/env"
	"api/modules/todos"
	"api/modules/users/cli"
	"api/modules/users/controllers"
	"api/modules/users/infra/hashpassword"
	"api/modules/users/infra/repositories"
	tokenjwt "api/modules/users/infra/token"
	"api/modules/users/middlewares"
	"api/modules/users/usecases"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// TODO: refatorar essa main em pequenas funcoes
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
	routerPublic := gin.Default()

	//cors
	routerPublic.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4000", "https://todo-frontend-from-golang.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	// auth secret key
	secretKey := env.TokenAuthSecretKey

	{
		// users routes public
		userRepository := repositories.NewUserRepository(db)
		hashPassword := hashpassword.NewHashPassword()
		userUsecase := usecases.NewUserUsecase(userRepository, hashPassword)
		userController := controllers.NewUserController(userUsecase)

		JWTMaker, err := tokenjwt.NewJWTMaker(secretKey)
		if err != nil {
			panic(err)
		}
		tokenManager := usecases.NewTokenManager(JWTMaker)
		loginUsecase := usecases.NewTokenLoginUsecase(userRepository, hashPassword, tokenManager)
		loginController := controllers.NewLoginController(loginUsecase)

		// create user genesis
		cliUserController := cli.NewUserController(userUsecase)
		err = cliUserController.AddUserGenesis()
		if err != nil {
			panic(err)
		}

		routerPublic.POST("/users", userController.CreateUser())

		// login routes private
		routerPublic.POST("/users/login", loginController.Login())

		{
			// users routes private
			JWTMaker, err := tokenjwt.NewJWTMaker(secretKey)
			if err != nil {
				panic(err)
			}
			tokenManager := usecases.NewTokenManager(JWTMaker)
			authMiddleware := middlewares.NewAuthorizationMiddleware(tokenManager)
			userRouterPrivate := routerPublic.Group("/")
			userRouterPrivate.Use(authMiddleware.Authorize())

			userRouterPrivate.PUT("/users/:id", userController.UpdateUser())
			userRouterPrivate.GET("/users/:id", userController.GetUser())
			userRouterPrivate.GET("/users", userController.GetAllUser())
			userRouterPrivate.DELETE("/users/:id", userController.DeleteUser())
			userRouterPrivate.POST("/users/change_password/:id", userController.ChangePassword())
			userRouterPrivate.PATCH("/users/photo/:id", userController.UpdatePhotoUser())
			userRouterPrivate.DELETE("/users/photo/:id", userController.DeletePhotoUser())
			userRouterPrivate.GET("/users/photo/:id", userController.GetPhotoUser())

		}
	}

	{
		// todos routes private
		JWTMaker, err := tokenjwt.NewJWTMaker(secretKey)
		if err != nil {
			panic(err)
		}
		tokenManager := usecases.NewTokenManager(JWTMaker)
		authMiddleware := middlewares.NewAuthorizationMiddleware(tokenManager)
		todoRouterPrivate := routerPublic.Group("/")
		todoRouterPrivate.Use(authMiddleware.Authorize())

		todoRepository := todos.NewTodoRepository(db)
		userRepository := repositories.NewUserRepository(db)
		hashPassword := hashpassword.NewHashPassword()
		userUsecase := usecases.NewUserUsecase(userRepository, hashPassword)
		todoUsecase := todos.NewTodoUsecase(todoRepository, userUsecase)
		controller := todos.NewTodoController(todoUsecase)
		todoRouterPrivate.POST("/todos", controller.CreateTodo())
		todoRouterPrivate.GET("/todos/:id", controller.GetTodo())
		todoRouterPrivate.GET("/todos", controller.GetAllTodos())
		todoRouterPrivate.PUT("/todos/:id", controller.UpdateTodo())
		todoRouterPrivate.DELETE("/todos/:id", controller.DeleteTodo())

		// todos image
		todoRouterPrivate.PATCH("/todos/image/:id", controller.UpdateImageTodo())
		todoRouterPrivate.GET("/todos/image/:id", controller.GetImageTodo())
		todoRouterPrivate.DELETE("/todos/image/:id", controller.DeleteImageTodo())

		// todo status
		todoRouterPrivate.POST("todos/status", controller.CreateStatusTodo())
		todoRouterPrivate.GET("todos/status/:id", controller.GetStatusTodo())
		todoRouterPrivate.GET("todos/status", controller.GetAllStatusTodo())
		todoRouterPrivate.PUT("todos/status/:id", controller.UpdateStatusTodo())
		todoRouterPrivate.DELETE("todos/status/:id", controller.DeleteStatusTodo())
	}

	routerPublic.Run()

}
