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

	"github.com/in-jun/go-structure-example/internal/shared/config"
	"github.com/in-jun/go-structure-example/internal/shared/database"
	"github.com/in-jun/go-structure-example/internal/shared/health"
	"github.com/in-jun/go-structure-example/internal/shared/logging"
	"github.com/in-jun/go-structure-example/internal/shared/middleware"
	sharedNats "github.com/in-jun/go-structure-example/internal/shared/nats"
	"github.com/in-jun/go-structure-example/internal/shared/observability"
	"github.com/in-jun/go-structure-example/internal/shared/outbox"
	"github.com/in-jun/go-structure-example/internal/shared/server"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"

	"github.com/in-jun/go-structure-example/internal/payment/application"
	"github.com/in-jun/go-structure-example/internal/payment/application/command"
	"github.com/in-jun/go-structure-example/internal/payment/application/query"
	"github.com/in-jun/go-structure-example/internal/payment/domain/service"
	"github.com/in-jun/go-structure-example/internal/payment/infrastructure/event"
	"github.com/in-jun/go-structure-example/internal/payment/infrastructure/gateway"
	paymentNats "github.com/in-jun/go-structure-example/internal/payment/infrastructure/nats"
	"github.com/in-jun/go-structure-example/internal/payment/infrastructure/pg"
	paymentHTTP "github.com/in-jun/go-structure-example/internal/payment/interfaces/http"
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
			port = "6064"
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
	logging.Init("payment-service")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	observability.InitMetrics()
	shutdownTracer, err := observability.InitTracer(ctx, "payment-service")
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

	pgDB, err := database.NewPostgres()
	if errors.Is(err, database.ErrMigrateOnly) {
		if err := pgDB.Close(); err != nil {
			slog.Error("failed to close db", "error", err)
		}
		os.Exit(0)
	}
	if err != nil {
		slog.Error("failed to connect to PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := pgDB.Close(); err != nil {
			slog.Error("failed to close db", "error", err)
		}
	}()

	nc, err := sharedNats.NewConnection()
	if err != nil {
		slog.Error("failed to connect to NATS", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := nc.Drain(); err != nil {
			slog.Warn("failed to drain NATS connection", "error", err)
		}
	}()

	dbGetter := transaction.NewDBGetter(pgDB)
	transactor := transaction.NewTransactor(pgDB)

	paymentRepo := pg.NewPaymentRepository(dbGetter)
	mockGW := gateway.NewMockGateway()
	processor := service.NewPaymentProcessor(mockGW)

	pgPublisher := event.NewPublisher(dbGetter)
	compositePublisher := event.NewCompositePublisher(pgPublisher, nc)

	createPaymentHandler := command.NewCreatePaymentHandler(paymentRepo, compositePublisher, transactor)
	confirmPaymentHandler := command.NewConfirmPaymentHandler(paymentRepo, processor, compositePublisher, transactor)
	refundPaymentHandler := command.NewRefundPaymentHandler(paymentRepo, processor, compositePublisher, transactor)
	getPaymentHandler := query.NewGetPaymentHandler(paymentRepo)

	consumer := paymentNats.NewConsumer(nc, createPaymentHandler, dbGetter, transactor)
	if err := consumer.Start(ctx); err != nil {
		slog.Error("failed to start NATS consumer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := consumer.Stop(); err != nil {
			slog.Warn("failed to stop consumer", "error", err)
		}
	}()

	relay := outbox.NewRelay(pgDB, nc, "payment")
	go relay.Start(ctx)

	svc := application.NewService(createPaymentHandler, confirmPaymentHandler, refundPaymentHandler, getPaymentHandler)

	var commands application.CommandUseCase = svc
	var queries application.QueryUseCase = svc

	handler := paymentHTTP.NewHandler(commands, queries)

	mux := server.NewRouter()

	stack := server.Chain(
		middleware.Recovery(),
		middleware.Timeout(30*time.Second),
		middleware.BodyLimit(1<<20),
		middleware.RequestID(),
		middleware.AccessLog(),
		middleware.CORS(config.AppConfig.CORSAllowOrigins),
		middleware.SecurityHeaders(),
		middleware.Tracing("payment-service"),
		middleware.Metrics("payment-service"),
	)

	mux.Handle("GET /metrics", observability.MetricsHandler())

	healthChecker := health.NewChecker(pgDB, nc).WithBuildInfo(Version, BuildTime, GitCommit)
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
		slog.Info("service starting", "service", "payment-service", "port", config.AppConfig.AppPort)
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
