package api

import (
	"encoding/json"
	"net/http"

	"GoCourse/logger"
	"GoCourse/store"
)

// writeJSON sends a JSON response with given status
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// CreateHandler adds a new task
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var newItem store.TodoItem
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	tasks := store.SendCommand(store.TaskCommand{
		Type:    store.CreateTask,
		NewTask: newItem,
	})
	logger.ContextLogger(r.Context()).Info("Created new task")
	writeJSON(w, http.StatusCreated, tasks[len(tasks)-1]) // return just the new one
}

// GetHandler returns all tasks
func GetHandler(w http.ResponseWriter, r *http.Request) {
	tasks := store.SendCommand(store.TaskCommand{
		Type: store.GetAllTasks,
	})

	logger.ContextLogger(r.Context()).Info("Returned tasks")
	writeJSON(w, http.StatusOK, tasks)
}

// Update Handler to update tasks
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	// define a struct that matches the expected json body from the client
	// index = which task to update
	// desc/status are optional updates
	var req struct {
		Index  int    `json:"index"`
		Desc   string `json:"desc"`
		Status string `json:"status"`
	}

	// read and decode the json request body into req, return error if necessary
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// send an update command to the store actor to ensure only one goroutine modifies the tasks at a time
	tasks := store.SendCommand(store.TaskCommand{
		Type:   store.UpdateTask,
		Index:  req.Index,
		Desc:   req.Desc,
		Status: req.Status,
	})

	// log that a task was updated (trace id comes from middleware) and return the updated task back to the client as json
	logger.ContextLogger(r.Context()).Info("updated a task")
	writeJSON(w, http.StatusOK, tasks[req.Index])
}

// handler to delete task
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// define a struct to read the index of the task to delete
	var req struct {
		Index int `json:"index"`
	}

	// decode the json request body, if the body is invalid, return a 400 error
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// send a delete command to the store actor - the actor safely removes the task and updates the file
	store.SendCommand(store.TaskCommand{
		Type:  store.DeleteaTask,
		Index: req.Index,
	})

	// log that a task was deleted and return a simple success response to the client
	logger.ContextLogger(r.Context()).Info("deleted a task")

	writeJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
}
