package shell

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
)

// Start Trace Endpoints.
//
// /trace/pprof
// /trace/vars
func (s *Shell) runTrace() {
	if os.Getenv("RF_TRACE_PORT") == "" {
		s.shell.Println(msgPrompt + "Error: RF_TRACE_PORT environment variable is not set.")
		return
	}

	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Mount("/trace", chiMiddleware.Profiler())

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", os.Getenv("RF_TRACE_PORT")),
		Handler: r,
	}

	go func() {
		s.shell.Println(msgPrompt + "Trace endpoint starting on port " + os.Getenv("RF_TRACE_PORT"))

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.shell.Println(msgPrompt+"HTTP server error:", err)
		}

		s.shell.Println(msgPrompt + "Trace endpoint stopped")
	}()

	go func() {
		<-s.ctx.Done()
		s.shell.Println(msgPrompt + "Trace endpoint shutting down...")
		if err := server.Shutdown(context.Background()); err != nil {
			s.shell.Println(msgPrompt+"Error shutting down server:", err)
			return
		}

		s.shell.Println(msgPrompt + "Trace endpoint shutdown gracefully.")
	}()

}
