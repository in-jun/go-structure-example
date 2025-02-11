package main

import (
	"log"

	"github.com/in-jun/go-structure-example/internal/app/auth"
	"github.com/in-jun/go-structure-example/internal/app/todo"
	"github.com/in-jun/go-structure-example/internal/app/user"
	"github.com/in-jun/go-structure-example/internal/pkg/config"
	"github.com/in-jun/go-structure-example/internal/pkg/db/mysql"
	"github.com/in-jun/go-structure-example/internal/pkg/db/redis"
	"github.com/in-jun/go-structure-example/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()

	redisClient, err := redis.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	mysqlDB, err := mysql.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer mysqlDB.Close()

	authRepo := redis.NewAuthRepository(redisClient)
	userRepo := mysql.NewUserRepository(mysqlDB)
	todoRepo := mysql.NewTodoRepository(mysqlDB)

	authService := auth.NewService(authRepo, userRepo)
	userService := user.NewService(userRepo)
	todoService := todo.NewService(todoRepo)

	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService)
	todoHandler := todo.NewHandler(todoService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())

	api := router.Group("/api/v1")
	{
		authHandler.RegisterRoutes(api)
		userHandler.RegisterRoutes(api)
		todoHandler.RegisterRoutes(api)
	}

	log.Printf("Starting server on port %s", config.AppConfig.AppPort)
	if err := router.Run(":" + config.AppConfig.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
