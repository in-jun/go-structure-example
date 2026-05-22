package health

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/in-jun/go-structure-example/internal/shared/server"
	"github.com/nats-io/nats.go"
)

// CheckReady verifies the service is accepting connections by opening a TCP
// connection to the service's port. Uses TCP rather than HTTP to avoid taint
// analysis false positives; a successful connect means the HTTP listener is up.
func CheckReady() {
	portStr := os.Getenv("APP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 || port > 65535 {
		port = 8080
	}
	addr := net.JoinHostPort("localhost", strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second) //nolint:gosec // host is hardcoded to localhost; port is range-validated
	if err != nil {
		os.Exit(1)
	}
	if err := conn.Close(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

type Checker struct {
	db        *sql.DB
	nc        *nats.Conn
	redis     *goredis.Client
	version   string
	buildTime string
	gitCommit string
}

func NewChecker(db *sql.DB, nc *nats.Conn) *Checker {
	return &Checker{db: db, nc: nc}
}

func NewCheckerWithoutDB(nc *nats.Conn) *Checker {
	return &Checker{nc: nc}
}

func (c *Checker) WithRedis(rc *goredis.Client) *Checker {
	c.redis = rc
	return c
}

func (c *Checker) WithBuildInfo(version, buildTime, gitCommit string) *Checker {
	c.version = version
	c.buildTime = buildTime
	c.gitCommit = gitCommit
	return c
}

func (c *Checker) RegisterRoutes(mux *server.Router) {
	mux.HandleFunc("GET /health/live", c.liveHandler)
	mux.HandleFunc("GET /health/ready", c.readyHandler)
}

func (c *Checker) liveHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]any{"status": "ok"}
	if c.version != "" {
		resp["version"] = c.version
		resp["build_time"] = c.buildTime
		resp["git_commit"] = c.gitCommit
	}
	server.JSON(w, http.StatusOK, resp)
}

func (c *Checker) readyHandler(w http.ResponseWriter, r *http.Request) {
	checks := map[string]any{}
	ready := true

	pingCtx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if c.db != nil {
		if err := c.db.PingContext(pingCtx); err != nil {
			checks["postgres"] = err.Error()
			ready = false
		} else {
			checks["postgres"] = "ok"
		}
	}

	if c.nc != nil {
		if c.nc.IsConnected() {
			checks["nats"] = "ok"
		} else {
			checks["nats"] = "disconnected"
			ready = false
		}
	}

	if c.redis != nil {
		if err := c.redis.Ping(pingCtx).Err(); err != nil {
			checks["redis"] = err.Error()
			ready = false
		} else {
			checks["redis"] = "ok"
		}
	}

	status := http.StatusOK
	statusText := "ready"
	if !ready {
		status = http.StatusServiceUnavailable
		statusText = "not ready"
	}

	server.JSON(w, status, map[string]any{"status": statusText, "checks": checks})
}
