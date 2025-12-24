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

// CreateHandler adds a new task
// func CreateHandler(w http.ResponseWriter, r *http.Request) {
// 	var newItem store.TodoItem
// 	err := json.NewDecoder(r.Body).Decode(&newItem)
// 	if err != nil {
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 		return
// 	}

// 	items, err := store.LoadTasks("tasks.json")
// 	if err != nil {
// 		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
// 		return
// 	}

// 	items = append(items, newItem)

// 	err = store.SaveTasks("tasks.json", items)
// 	if err != nil {
// 		http.Error(w, "Failed to save tasks", http.StatusInternalServerError)
// 		return
// 	}
// 	logger.ContextLogger(r.Context()).Info("Created new task")
// 	writeJSON(w, http.StatusCreated, newItem)
// }

// GetHandler returns all tasks
func GetHandler(w http.ResponseWriter, r *http.Request) {
	tasks := store.SendCommand(store.TaskCommand{
		Type: store.GetAllTasks,
	})

	logger.ContextLogger(r.Context()).Info("Returned tasks")
	writeJSON(w, http.StatusOK, tasks)
}

// GetHandler lists all tasks
// func GetHandler(w http.ResponseWriter, r *http.Request) {
// 	items, err := store.LoadTasks("tasks.json")
// 	if err != nil {
// 		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
// 		return
// 	}
// 	logger.ContextLogger(r.Context()).Info("Returned tasks")
// 	writeJSON(w, http.StatusOK, items)
// }

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

	// send an update command to the store actor
	// this ensures only one goroutine modifies the tasks at a time
	tasks := store.SendCommand(store.TaskCommand{
		Type:   store.UpdateTask,
		Index:  req.Index,
		Desc:   req.Desc,
		Status: req.Status,
	})

	// log that a task was updated (trace id comes from middleware)
	logger.ContextLogger(r.Context()).Info("updated a task")

	// return the updated task back to the client as json
	writeJSON(w, http.StatusOK, tasks[req.Index])
}

// func UpdateHandler(w http.ResponseWriter, r *http.Request) {
// 	var req struct {
// 		Index  int    `json:"index"`
// 		Desc   string `json:"desc"`
// 		Status string `json:"status"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	items, err := store.LoadTasks("tasks.json")
// 	if err != nil || req.Index < 0 || req.Index >= len(items) {
// 		http.Error(w, "Invalid index", http.StatusBadRequest)
// 		return
// 	}

// 	if req.Desc != "" {
// 		items[req.Index].Description = req.Desc
// 	}
// 	if req.Status != "" {
// 		items[req.Index].Status = req.Status
// 	}

// 	if err := store.SaveTasks("tasks.json", items); err != nil {
// 		http.Error(w, "Failed to save", http.StatusInternalServerError)
// 		return
// 	}
// 	logger.ContextLogger(r.Context()).Info("Updated a task")
// 	w.WriteHeader(http.StatusNoContent)
// }

// handler to delete task
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// define a struct to read the index of the task to delete
	var req struct {
		Index int `json:"index"`
	}

	// decode the json request body
	// if the body is invalid, return a 400 error
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// send a delete command to the store actor
	// the actor safely removes the task and updates the file
	store.SendCommand(store.TaskCommand{
		Type:  store.DeleteaTask,
		Index: req.Index,
	})

	// log that a task was deleted
	logger.ContextLogger(r.Context()).Info("deleted a task")

	// return a simple success response to the client
	writeJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
}

// func DeleteHandler(w http.ResponseWriter, r *http.Request) {
// 	var req struct {
// 		Index int `json:"index"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	items, err := store.LoadTasks("tasks.json")
// 	if err != nil || req.Index < 0 || req.Index >= len(items) {
// 		http.Error(w, "Invalid index", http.StatusBadRequest)
// 		return
// 	}

// 	items = append(items[:req.Index], items[req.Index+1:]...)

// 	if err := store.SaveTasks("tasks.json", items); err != nil {
// 		http.Error(w, "Failed to save", http.StatusInternalServerError)
// 		return
// 	}
// 	logger.ContextLogger(r.Context()).Info("deleted a task")
// 	w.WriteHeader(http.StatusNoContent)
// }
