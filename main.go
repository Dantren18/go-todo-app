package main

import (
	"GoCourse/api"
	"GoCourse/logger"
	"GoCourse/store"
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	store.StartStoreActor()

	mux := http.NewServeMux()
	mux.HandleFunc("/create", api.CreateHandler)
	mux.HandleFunc("/get", api.GetHandler)
	mux.HandleFunc("/update", api.UpdateHandler)
	mux.HandleFunc("/delete", api.DeleteHandler)

	// Serve static about.html from web/static/
	static := http.FileServer(http.Dir("web/static"))
	mux.Handle("/about/", http.StripPrefix("/about/", static))

	// Serve dynamic list of tasks at /list
	tmpl := template.Must(template.ParseFiles("web/templates/list.html"))
	mux.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := store.LoadTasks("tasks.json")
		if err != nil {

		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, tasks)
	})

	// Wrap mux with TraceID middleware
	handler := logger.TraceIDMiddleware(mux)

	// Create an http.Server so we can shut it down cleanly later
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Run the server in a goroutine. ListenAndServe blocks, so running it
	// here lets the main goroutine wait for OS signals below.
	go func() {
		slog.Info("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Only log unexpected errors (ErrServerClosed is returned by Shutdown)
			slog.Error("server error", "err", err)
		}
	}()

	// Set up channel on which to receive OS signals. We use a small buffer
	// so the signal isn't missed if nobody is ready to receive immediately.
	quit := make(chan os.Signal, 1)
	// Notify the channel when an interrupt (Ctrl+C) or SIGTERM is received.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block here until we receive a signal.
	<-quit

	slog.Info("Shutdown signal received, shutting down server...")

	// Create a context with timeout for the graceful shutdown. This gives
	// inflight requests up to 5 seconds to finish before we exit.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "err", err)
	} else {
		slog.Info("server stopped gracefully")
	}
}
