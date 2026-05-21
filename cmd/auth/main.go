package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "go.uber.org/automaxprocs"

	"github.com/in-jun/go-structure-example/internal/auth/application"
	"github.com/in-jun/go-structure-example/internal/auth/application/command"
	"github.com/in-jun/go-structure-example/internal/auth/application/query"
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
	"github.com/in-jun/go-structure-example/internal/shared/observability"
	"github.com/in-jun/go-structure-example/internal/shared/server"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"

	goredis "github.com/go-redis/redis/v8"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost:"+os.Getenv("APP_PORT")+"/health/ready", nil)
		if err != nil {
			os.Exit(1)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			os.Exit(1)
		}
		status := resp.StatusCode
		if err := resp.Body.Close(); err != nil {
			slog.Warn("failed to close healthcheck response body", "error", err)
		}
		if status != http.StatusOK {
			os.Exit(1)
		}
		os.Exit(0)
	}

	go func() {
		port := os.Getenv("PPROF_PORT")
		if port == "" {
			port = "6061"
		}
		pprofSrv := &http.Server{
			Addr:         "localhost:" + port,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		}
		if err := pprofSrv.ListenAndServe(); err != nil {
			slog.Warn("pprof server stopped", "error", err)
		}
	}()

	config.Load()
	logging.Init("auth-service")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	observability.InitMetrics()
	shutdownTracer, err := observability.InitTracer(ctx, "auth-service")
	if err != nil {
		slog.Warn("failed to init tracer", "error", err)
	}
	if shutdownTracer != nil {
		defer func() {
			if err := shutdownTracer(context.Background()); err != nil {
				slog.Warn("failed to shutdown tracer", "error", err)
			}
		}()
	}

	db, err := database.NewPostgres()
	if errors.Is(err, database.ErrMigrateOnly) {
		if err := db.Close(); err != nil {
			slog.Error("failed to close db", "error", err)
		}
		os.Exit(0)
	}
	if err != nil {
		slog.Error("failed to connect to PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close db", "error", err)
		}
	}()

	redisClient := goredis.NewClient(&goredis.Options{
		Addr: config.AppConfig.RedisURL,
	})
	defer func() {
		if err := redisClient.Close(); err != nil {
			slog.Warn("failed to close Redis client", "error", err)
		}
	}()

	dbGetter := transaction.NewDBGetter(db)
	transactor := transaction.NewTransactor(db)

	tokenGen, err := authjwt.NewProvider(
		config.AppConfig.JWTSecret,
		config.AppConfig.JWTAccessExpiry,
		config.AppConfig.JWTRefreshExpiry,
	)
	if err != nil {
		slog.Error("failed to create JWT provider", "error", err)
		os.Exit(1)
	}

	hasher := crypto.NewBcryptPasswordHasher()
	userRepo := authpg.NewUserRepository(dbGetter)
	tokenRepo := authredis.NewTokenRepository(redisClient)

	svc := application.NewService(
		command.NewRegisterHandler(userRepo, hasher, transactor),
		command.NewLoginHandler(userRepo, tokenRepo, tokenGen, hasher),
		command.NewRefreshHandler(tokenRepo, tokenGen),
		command.NewLogoutHandler(tokenRepo, tokenGen),
		command.NewLogoutAllHandler(tokenRepo, tokenGen),
		query.NewValidateHandler(tokenRepo, tokenGen),
	)

	var commands application.CommandUseCase = svc
	var queries application.QueryUseCase = svc

	tokenValidator := middleware.TokenValidator(func(ctx context.Context, tokenString string) (*middleware.ValidateTokenResult, error) {
		result, err := queries.ValidateToken(ctx, query.Validate{TokenString: tokenString})
		if err != nil {
			return nil, err
		}
		return &middleware.ValidateTokenResult{UserID: result.UserID, JTI: result.JTI, IssuedAt: result.IssuedAt}, nil
	})

	handler := authhttp.NewHandler(commands, queries, tokenValidator)

	mux := server.NewRouter()

	stack := server.Chain(
		middleware.Recovery(),
		middleware.Timeout(30*time.Second),
		middleware.BodyLimit(1<<20),
		middleware.RequestID(),
		middleware.AccessLog(),
		middleware.CORS(config.AppConfig.CORSAllowOrigins),
		middleware.SecurityHeaders(),
		middleware.Tracing("auth-service"),
		middleware.Metrics("auth-service"),
	)

	mux.Handle("GET /metrics", observability.MetricsHandler())

	healthChecker := health.NewChecker(db, nil).WithRedis(redisClient).WithBuildInfo(Version, BuildTime, GitCommit)
	healthChecker.RegisterRoutes(mux)

	handler.RegisterRoutes(mux, stack)

	srv := &http.Server{
		Addr:         ":" + config.AppConfig.AppPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("service starting", "service", "auth-service", "port", config.AppConfig.AppPort)
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

	slog.Info("service stopped")
}
