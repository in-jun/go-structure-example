package observability

import (
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"time"
)

// StartPprofServer starts a pprof profiling server on localhost only.
// It registers handlers on an explicit ServeMux (not DefaultServeMux) to
// prevent accidental exposure through other servers that use DefaultServeMux.
func StartPprofServer() {
	port := os.Getenv("PPROF_PORT")
	if port == "" {
		port = "6062"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := &http.Server{
		Addr:         "localhost:" + port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		slog.Warn("pprof server stopped", "error", err)
	}
}
