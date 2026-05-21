package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	_ "net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "go.uber.org/automaxprocs"

	goredis "github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/in-jun/go-structure-example/internal/shared/config"
	"github.com/in-jun/go-structure-example/internal/shared/health"
	"github.com/in-jun/go-structure-example/internal/shared/logging"
	"github.com/in-jun/go-structure-example/internal/shared/middleware"
	"github.com/in-jun/go-structure-example/internal/shared/observability"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

type serviceProxy struct {
	proxy *httputil.ReverseProxy
	name  string
}

func newServiceProxy(name string, rawURL string, timeout time.Duration, maxRetries int) *serviceProxy {
	target, _ := url.Parse(rawURL)

	base := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		ResponseHeaderTimeout: timeout,
		TLSHandshakeTimeout:   5 * time.Second,
	}

	retry := middleware.NewRetryTransport(base, maxRetries)
	cb := middleware.NewCircuitBreakerTransport(retry, name)

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = cb
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		slog.Error("proxy error", "service", name, "path", r.URL.Path, "error", err)
		server.Error(w, http.StatusBadGateway, "Bad Gateway")
	}

	return &serviceProxy{proxy: proxy, name: name}
}

func isTokenBlacklisted(ctx context.Context, redisClient *goredis.Client, jti string, userID string, issuedAt int64) bool {
	pipe := redisClient.Pipeline()

	var jtiCmd *goredis.IntCmd
	if jti != "" {
		jtiCmd = pipe.Exists(ctx, "token_blacklist:"+jti)
	}
	revokedCmd := pipe.Get(ctx, "token_revoked_at:"+userID)

	if _, err := pipe.Exec(ctx); err != nil {
		slog.Warn("failed to exec redis pipeline for token blacklist check", "error", err)
	}

	if jtiCmd != nil {
		if val, err := jtiCmd.Result(); err == nil && val > 0 {
			return true
		}
	}

	if revokedAtStr, err := revokedCmd.Result(); err == nil {
		if revokedAt, err := strconv.ParseInt(revokedAtStr, 10, 64); err == nil {
			if issuedAt <= revokedAt {
				return true
			}
		}
	}

	return false
}

type jwtClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

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
			port = "6065"
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
	logging.Init("gateway-service")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	observability.InitMetrics()
	shutdownTracer, err := observability.InitTracer(ctx, "gateway-service")
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

	redisClient := goredis.NewClient(&goredis.Options{
		Addr: config.AppConfig.RedisURL,
	})
	defer func() {
		if err := redisClient.Close(); err != nil {
			slog.Warn("failed to close Redis client", "error", err)
		}
	}()

	secretKey := []byte(config.AppConfig.JWTSecret)
	tokenValidator := middleware.TokenValidator(func(ctx context.Context, tokenString string) (*middleware.ValidateTokenResult, error) {
		token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})
		if err != nil {
			return nil, err
		}
		c, ok := token.Claims.(*jwtClaims)
		if !ok || !token.Valid {
			return nil, jwt.ErrSignatureInvalid
		}

		var issuedAt int64
		if c.IssuedAt != nil {
			issuedAt = c.IssuedAt.Unix()
		}

		if isTokenBlacklisted(ctx, redisClient, c.ID, c.UserID, issuedAt) {
			return nil, fmt.Errorf("token has been revoked")
		}

		return &middleware.ValidateTokenResult{
			UserID:   c.UserID,
			JTI:      c.ID,
			IssuedAt: issuedAt,
		}, nil
	})

	authSvc := newServiceProxy("auth", config.AppConfig.AuthServiceURL, 5*time.Second, 2)
	auctionSvc := newServiceProxy("auction", config.AppConfig.AuctionServiceURL, 10*time.Second, 3)
	bidSvc := newServiceProxy("bid", config.AppConfig.BidServiceURL, 10*time.Second, 3)
	paymentSvc := newServiceProxy("payment", config.AppConfig.PaymentServiceURL, 30*time.Second, 2)

	mux := server.NewRouter()

	stack := server.Chain(
		middleware.Recovery(),
		middleware.Timeout(30*time.Second),
		middleware.BodyLimit(1<<20),
		middleware.RequestID(),
		middleware.AccessLog(),
		middleware.CORS(config.AppConfig.CORSAllowOrigins),
		middleware.SecurityHeaders(),
		middleware.RateLimit(redisClient, config.AppConfig.RateLimitRPS, config.AppConfig.RateLimitBurst),
		middleware.Tracing("gateway-service"),
		middleware.Metrics("gateway-service"),
	)

	mux.Handle("GET /metrics", observability.MetricsHandler())

	healthChecker := health.NewCheckerWithoutDB(nil).WithRedis(redisClient).WithBuildInfo(Version, BuildTime, GitCommit)
	healthChecker.RegisterRoutes(mux)

	publicProxy := func(sp *serviceProxy) http.Handler {
		return stack(sp.proxy)
	}

	authMw := middleware.Auth(tokenValidator)
	idempotencyMw := middleware.Idempotency(redisClient)
	injectUserID := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set("X-User-ID", server.UserID(r))
			r.Header.Set("X-Token-JTI", server.TokenJTI(r))
			r.Header.Set("X-Token-IAT", strconv.FormatInt(server.TokenIssuedAt(r), 10))
			next.ServeHTTP(w, r)
		})
	}

	authedProxy := func(sp *serviceProxy) http.Handler {
		return stack(authMw(idempotencyMw(injectUserID(sp.proxy))))
	}

	authedNoIdempotency := func(sp *serviceProxy) http.Handler {
		return stack(authMw(injectUserID(sp.proxy)))
	}

	// Auth routes
	mux.Handle("POST /api/v1/auth/register", publicProxy(authSvc))
	mux.Handle("POST /api/v1/auth/login", publicProxy(authSvc))
	mux.Handle("POST /api/v1/auth/refresh", publicProxy(authSvc))
	mux.Handle("POST /api/v1/auth/logout", authedNoIdempotency(authSvc))
	mux.Handle("POST /api/v1/auth/logout/all", authedNoIdempotency(authSvc))

	// Auction routes (authed write, public read)
	mux.Handle("POST /api/v1/auctions", authedProxy(auctionSvc))
	mux.Handle("POST /api/v1/auctions/{id}/open", authedProxy(auctionSvc))
	mux.Handle("POST /api/v1/auctions/{id}/close", authedProxy(auctionSvc))
	mux.Handle("POST /api/v1/auctions/{id}/cancel", authedProxy(auctionSvc))
	mux.Handle("GET /api/v1/auctions", publicProxy(auctionSvc))
	mux.Handle("GET /api/v1/auctions/{id}", publicProxy(auctionSvc))
	mux.Handle("GET /api/v1/auctions/{id}/events", publicProxy(auctionSvc))

	// Bid routes
	mux.Handle("POST /api/v1/auctions/{id}/bids", authedProxy(bidSvc))
	mux.Handle("GET /api/v1/auctions/{id}/bids", publicProxy(bidSvc))
	mux.Handle("GET /api/v1/auctions/{id}/bids/highest", publicProxy(bidSvc))

	// Payment routes
	mux.Handle("POST /api/v1/payments/{id}/confirm", authedProxy(paymentSvc))
	mux.Handle("POST /api/v1/payments/{id}/refund", authedProxy(paymentSvc))
	mux.Handle("GET /api/v1/payments/{id}", authedNoIdempotency(paymentSvc))

	srv := &http.Server{
		Addr:         ":" + config.AppConfig.AppPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("service starting", "service", "gateway", "port", config.AppConfig.AppPort)
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
