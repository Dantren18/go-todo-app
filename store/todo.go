package store

import (
	"encoding/json"
	"errors"
	"os"
)

// TodoItem defines a to‑do task with description and status
type TodoItem struct {
	Description string `json:"description"`
	Status      string `json:"status"`
}

// LoadTasks loads to‑do items from a JSON file
func LoadTasks(filename string) ([]TodoItem, error) {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []TodoItem{}, nil
		}
		return nil, err
	}

	var items []TodoItem
	if err := json.Unmarshal(fileData, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// SaveTasks saves to‑do items to a JSON file
func SaveTasks(filename string, items []TodoItem) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// UpdateTaskDescription updates the description of an item at index
func UpdateTaskDescription(items []TodoItem, index int, newDesc string) ([]TodoItem, error) {
	if index < 0 || index >= len(items) {
		return items, errors.New("update index out of range")
	}
	items[index].Description = newDesc
	return items, nil
}

// UpdateTaskStatus updates the status of an item at index
func UpdateTaskStatus(items []TodoItem, index int, newStatus string) ([]TodoItem, error) {
	if index < 0 || index >= len(items) {
		return items, errors.New("status update index out of range")
	}
	// require exact match of allowed statuses
	if newStatus != "Not started" && newStatus != "Started" && newStatus != "Completed" {
		return items, errors.New("Status must be either \"Not started\", \"Started\", or \"Completed\"")
	}
	items[index].Status = newStatus
	return items, nil
}

// DeleteTask removes an item at index
func DeleteTask(items []TodoItem, index int) ([]TodoItem, error) {
	if index < 0 || index >= len(items) {
		return items, errors.New("delete index out of range")
	}
	items = append(items[:index], items[index+1:]...)
	return items, nil
}
