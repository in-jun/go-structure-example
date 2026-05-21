package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authapp "github.com/in-jun/go-structure-example/internal/auth/application"
	authcmd "github.com/in-jun/go-structure-example/internal/auth/application/command"
	authqry "github.com/in-jun/go-structure-example/internal/auth/application/query"
	authjwt "github.com/in-jun/go-structure-example/internal/auth/infrastructure/jwt"
	authpg "github.com/in-jun/go-structure-example/internal/auth/infrastructure/pg"
	authredis "github.com/in-jun/go-structure-example/internal/auth/infrastructure/redis"
	authhttp "github.com/in-jun/go-structure-example/internal/auth/interfaces/http"
	"github.com/in-jun/go-structure-example/internal/shared/config"
	"github.com/in-jun/go-structure-example/internal/shared/crypto"
	"github.com/in-jun/go-structure-example/internal/shared/database"
	"github.com/in-jun/go-structure-example/internal/shared/health"
	"github.com/in-jun/go-structure-example/internal/shared/logging"
	"github.com/in-jun/go-structure-example/internal/shared/middleware"
	"github.com/in-jun/go-structure-example/internal/shared/server"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
	todoapp "github.com/in-jun/go-structure-example/internal/todo/application"
	todocmd "github.com/in-jun/go-structure-example/internal/todo/application/command"
	todoqry "github.com/in-jun/go-structure-example/internal/todo/application/query"
	todopg "github.com/in-jun/go-structure-example/internal/todo/infrastructure/pg"
	todohttp "github.com/in-jun/go-structure-example/internal/todo/interfaces/http"
	userapp "github.com/in-jun/go-structure-example/internal/user/application"
	usercmd "github.com/in-jun/go-structure-example/internal/user/application/command"
	userqry "github.com/in-jun/go-structure-example/internal/user/application/query"
	userpg "github.com/in-jun/go-structure-example/internal/user/infrastructure/pg"
	userhttp "github.com/in-jun/go-structure-example/internal/user/interfaces/http"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		resp, err := http.Get("http://localhost:" + os.Getenv("APP_PORT") + "/health/ready")
		if err != nil || resp.StatusCode != 200 {
			os.Exit(1)
		}
		os.Exit(0)
	}

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

	db, err := database.NewPostgres()
	if err != nil {
		slog.Error("failed to connect to PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	tokenGen, err := authjwt.NewProvider(
		config.AppConfig.JWTSecret,
		config.AppConfig.JWTAccessExpiry,
		config.AppConfig.JWTRefreshExpiry,
	)
	if err != nil {
		slog.Error("failed to create token generator", "error", err)
		os.Exit(1)
	}

	hasher := crypto.NewBcryptPasswordHasher()

	dbGetter := transaction.NewDBGetter(db)
	transactor := transaction.NewTransactor(db)

	tokenRepo := authredis.NewTokenRepository(redisClient)
	authUserRepo := authpg.NewUserRepository(dbGetter)
	userRepo := userpg.NewUserRepository(dbGetter)
	todoRepo := todopg.NewTodoRepository(dbGetter)

	authService := authapp.NewService(
		authcmd.NewRegisterHandler(authUserRepo, hasher, transactor),
		authcmd.NewLoginHandler(authUserRepo, tokenRepo, tokenGen, hasher),
		authcmd.NewRefreshHandler(tokenRepo, tokenGen),
		authcmd.NewLogoutHandler(tokenRepo, tokenGen),
		authcmd.NewLogoutAllHandler(tokenRepo, tokenGen),
		authqry.NewValidateHandler(tokenRepo, tokenGen),
	)

	var authQueries authapp.QueryUseCase = authService

	validateToken := middleware.TokenValidator(func(ctx context.Context, tokenString string) (*middleware.ValidateTokenResult, error) {
		result, err := authQueries.ValidateToken(ctx, authqry.Validate{TokenString: tokenString})
		if err != nil {
			return nil, err
		}
		return &middleware.ValidateTokenResult{UserID: result.UserID, JTI: result.JTI, IssuedAt: result.IssuedAt}, nil
	})

	userService := userapp.NewService(
		usercmd.NewUpdateProfileHandler(userRepo),
		usercmd.NewUpdatePasswordHandler(userRepo, hasher),
		usercmd.NewDeleteHandler(userRepo),
		userqry.NewGetHandler(userRepo),
	)

	todoService := todoapp.NewService(
		todocmd.NewCreateHandler(todoRepo),
		todocmd.NewUpdateHandler(todoRepo),
		todocmd.NewUpdateStatusHandler(todoRepo),
		todocmd.NewDeleteHandler(todoRepo),
		todoqry.NewGetHandler(todoRepo),
		todoqry.NewListHandler(todoRepo),
	)

	var userCommands userapp.CommandUseCase = userService
	var userQueries userapp.QueryUseCase = userService

	var todoCommands todoapp.CommandUseCase = todoService
	var todoQueries todoapp.QueryUseCase = todoService

	authHandler := authhttp.NewHandler(authService, authQueries, validateToken)
	userHandler := userhttp.NewHandler(userCommands, userQueries, validateToken)
	todoHandler := todohttp.NewHandler(todoCommands, todoQueries, validateToken)

	mux := server.NewRouter()

	stack := server.Chain(
		middleware.Recovery(),
		middleware.RequestID(),
		middleware.Timeout(30*time.Second),
		middleware.BodyLimit(1<<20),
		middleware.SecurityHeaders(),
		middleware.AccessLog(),
		middleware.CORS(config.AppConfig.CORSAllowOrigins),
		middleware.RateLimit(redisClient, config.AppConfig.RateLimitRPS, config.AppConfig.RateLimitBurst),
	)

	healthChecker := health.NewChecker(db, nil).WithRedis(redisClient).WithBuildInfo(Version, BuildTime, GitCommit)
	healthChecker.RegisterRoutes(mux)

	authHandler.RegisterRoutes(mux, stack)
	userHandler.RegisterRoutes(mux, stack)
	todoHandler.RegisterRoutes(mux, stack)

	srv := &http.Server{
		Addr:         ":" + config.AppConfig.AppPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
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
