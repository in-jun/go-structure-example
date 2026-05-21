package health

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

type Checker struct {
	db        *sql.DB
	redis     *goredis.Client
	version   string
	buildTime string
	gitCommit string
}

func NewChecker(db *sql.DB) *Checker {
	return &Checker{db: db}
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

	if c.db != nil {
		tctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := c.db.PingContext(tctx); err != nil {
			checks["postgres"] = "unhealthy: " + err.Error()
			ready = false
		} else {
			checks["postgres"] = "healthy"
		}
	}

	if c.redis != nil {
		tctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := c.redis.Ping(tctx).Err(); err != nil {
			checks["redis"] = "unhealthy: " + err.Error()
			ready = false
		} else {
			checks["redis"] = "healthy"
		}
	}

	status := http.StatusOK
	statusStr := "ok"
	if !ready {
		status = http.StatusServiceUnavailable
		statusStr = "degraded"
	}
	server.JSON(w, status, map[string]any{"status": statusStr, "checks": checks})
}
