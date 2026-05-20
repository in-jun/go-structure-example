package health

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v8"
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

func (c *Checker) RegisterRoutes(r *gin.Engine) {
	r.GET("/health/live", c.liveHandler)
	r.GET("/health/ready", c.readyHandler)
}

func (c *Checker) liveHandler(ctx *gin.Context) {
	resp := gin.H{"status": "ok"}
	if c.version != "" {
		resp["version"] = c.version
		resp["build_time"] = c.buildTime
		resp["git_commit"] = c.gitCommit
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *Checker) readyHandler(ctx *gin.Context) {
	checks := map[string]any{}
	ready := true

	if c.db != nil {
		tctx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
		defer cancel()
		if err := c.db.PingContext(tctx); err != nil {
			checks["mysql"] = "unhealthy: " + err.Error()
			ready = false
		} else {
			checks["mysql"] = "healthy"
		}
	}

	if c.redis != nil {
		tctx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
		defer cancel()
		if err := c.redis.Ping(tctx).Err(); err != nil {
			checks["redis"] = "unhealthy: " + err.Error()
			ready = false
		} else {
			checks["redis"] = "healthy"
		}
	}

	status := http.StatusOK
	if !ready {
		status = http.StatusServiceUnavailable
	}
	ctx.JSON(status, gin.H{"status": map[bool]string{true: "ok", false: "degraded"}[ready], "checks": checks})
}
