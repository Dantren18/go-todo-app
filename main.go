package main

import (
	"GoCourse/api"
	"GoCourse/logger"
	"GoCourse/store"
	"html/template"
	"log"
	"log/slog"
	"net/http"
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

	// Log and start server on port 8080
	slog.Info("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
