package middleware

import (
	stderrors "errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

func Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					slog.Error("panic recovered", "panic", rec, "stack", string(debug.Stack()))
					server.Error(w, http.StatusInternalServerError, "Internal Server Error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func HandleError(w http.ResponseWriter, err error) {
	var customErr errors.CustomError
	if stderrors.As(err, &customErr) {
		if customErr.Status >= 500 {
			slog.Error("internal error", "message", customErr.Message)
			customErr.Message = "Internal Server Error"
		}
		server.JSON(w, customErr.Status, customErr)
		return
	}
	slog.Error("unhandled error", "error", err)
	server.JSON(w, http.StatusInternalServerError, errors.Internal("Internal Server Error"))
}
