package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	authapp "github.com/in-jun/go-structure-example/internal/auth/application"
	authcmd "github.com/in-jun/go-structure-example/internal/auth/application/command"
	authqry "github.com/in-jun/go-structure-example/internal/auth/application/query"
	authjwt "github.com/in-jun/go-structure-example/internal/auth/infrastructure/jwt"
	authmysql "github.com/in-jun/go-structure-example/internal/auth/infrastructure/mysql"
	authredis "github.com/in-jun/go-structure-example/internal/auth/infrastructure/redis"
	authhttp "github.com/in-jun/go-structure-example/internal/auth/interfaces/http"
	"github.com/in-jun/go-structure-example/internal/shared/config"
	"github.com/in-jun/go-structure-example/internal/shared/crypto"
	"github.com/in-jun/go-structure-example/internal/shared/database"
	"github.com/in-jun/go-structure-example/internal/shared/logging"
	"github.com/in-jun/go-structure-example/internal/shared/middleware"
	todoapp "github.com/in-jun/go-structure-example/internal/todo/application"
	todocmd "github.com/in-jun/go-structure-example/internal/todo/application/command"
	todoqry "github.com/in-jun/go-structure-example/internal/todo/application/query"
	todomysql "github.com/in-jun/go-structure-example/internal/todo/infrastructure/mysql"
	todohttp "github.com/in-jun/go-structure-example/internal/todo/interfaces/http"
	userapp "github.com/in-jun/go-structure-example/internal/user/application"
	usercmd "github.com/in-jun/go-structure-example/internal/user/application/command"
	userqry "github.com/in-jun/go-structure-example/internal/user/application/query"
	usermysql "github.com/in-jun/go-structure-example/internal/user/infrastructure/mysql"
	userhttp "github.com/in-jun/go-structure-example/internal/user/interfaces/http"
)

func main() {
	config.Load()
	logging.Init("go-structure-example")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	redisClient, err := database.NewRedis()
	if err != nil {
		slog.Error("failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	mysqlDB, err := database.NewMySQL()
	if err != nil {
		slog.Error("failed to connect to MySQL", "error", err)
		os.Exit(1)
	}
	defer mysqlDB.Close()

	tokenGen, err := authjwt.NewProvider(
		config.AppConfig.JWTSecret,
		config.AppConfig.JWTAccessExpiry,
		config.AppConfig.JWTRefreshExpiry,
	)
	if err != nil {
		slog.Error("failed to create token generator", "error", err)
		os.Exit(1)
	}

	hasher := crypto.NewBcryptHasher()

	tokenRepo := authredis.NewTokenRepository(redisClient)
	authUserRepo := authmysql.NewUserRepository(mysqlDB)
	userRepo := usermysql.NewUserRepository(mysqlDB)
	todoRepo := todomysql.NewTodoRepository(mysqlDB)

	authService := authapp.NewService(
		authcmd.NewRegisterHandler(authUserRepo, hasher),
		authcmd.NewLoginHandler(authUserRepo, tokenRepo, tokenGen, hasher),
		authcmd.NewRefreshHandler(tokenRepo, tokenGen),
		authcmd.NewLogoutHandler(tokenRepo, tokenGen),
		authcmd.NewLogoutAllHandler(tokenRepo, tokenGen),
		authqry.NewValidateHandler(tokenRepo, tokenGen),
	)

	var authQueries authapp.QueryUseCase = authService

	validateToken := middleware.TokenValidator(func(ctx context.Context, tokenString string) (*middleware.TokenValidateResult, error) {
		result, err := authQueries.ValidateToken(ctx, authqry.Validate{TokenString: tokenString})
		if err != nil {
			return nil, err
		}
		return &middleware.TokenValidateResult{UserID: result.UserID, JTI: result.JTI}, nil
	})

	userService := userapp.NewService(
		usercmd.NewUpdateProfileHandler(userRepo),
		usercmd.NewUpdatePasswordHandler(userRepo, hasher),
		usercmd.NewDeleteHandler(userRepo),
		userqry.NewGetUserHandler(userRepo),
	)

	todoService := todoapp.NewService(
		todocmd.NewCreateHandler(todoRepo),
		todocmd.NewUpdateHandler(todoRepo),
		todocmd.NewUpdateStatusHandler(todoRepo),
		todocmd.NewDeleteHandler(todoRepo),
		todoqry.NewGetTodoHandler(todoRepo),
		todoqry.NewListTodosHandler(todoRepo),
	)

	authHandler := authhttp.NewHandler(authService, authQueries, validateToken)
	userHandler := userhttp.NewHandler(userService, userService, validateToken)
	todoHandler := todohttp.NewHandler(todoService, todoService, validateToken)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.AccessLog())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimit(redisClient, config.AppConfig.RateLimitBurst))
	router.Use(middleware.ErrorHandler())

	api := router.Group("/api/v1")
	{
		authHandler.RegisterRoutes(api)
		userHandler.RegisterRoutes(api)
		todoHandler.RegisterRoutes(api)
	}

	srv := &http.Server{
		Addr:    ":" + config.AppConfig.AppPort,
		Handler: router,
	}

	go func() {
		slog.Info("server starting", "port", config.AppConfig.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.AppConfig.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server stopped")
}
