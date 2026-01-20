package store

import (
	"encoding/json"
	"os"
	"testing"
)

// Test loading tasks from a JSON file
func TestLoadTasks_JSON(t *testing.T) {
	filename := "test_tasks.json"
	sample := `[{"description":"Task A","status":"Not started"},{"description":"Task B","status":"Started"}]`
	err := os.WriteFile(filename, []byte(sample), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	defer os.Remove(filename)

	items, err := LoadTasks(filename)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}
	// Using %+v to print the entire struct for better error messages
	if items[0].Desc != "Task A" || items[0].Status != "Not started" {
		t.Errorf("First item incorrect: got %+v", items[0])
	}
	if items[1].Desc != "Task B" || items[1].Status != "Started" {
		t.Errorf("Second item incorrect: got %+v", items[1])
	}
}

// Test saving tasks to JSON file
func TestSaveTasks_JSON(t *testing.T) {
	filename := "test_save.json"
	defer os.Remove(filename)

	original := []TodoItem{
		{Desc: "Task X", Status: "Not started"},
		{Desc: "Task Y", Status: "Completed"},
	}
	err := SaveTasks(filename, original)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Could not read file: %v", err)
	}

	var loaded []TodoItem
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if len(loaded) != 2 || loaded[1].Status != "Completed" || loaded[0].Desc != "Task X" {
		t.Errorf("Expected tasks to match original, got %+v", loaded)
	}
}

// Test updating a task description
func TestUpdateTaskDescription(t *testing.T) {
	items := []TodoItem{{Desc: "Old", Status: "Not started"}}
	updated, err := UpdateTaskDescription(items, 0, "New")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if updated[0].Desc != "New" {
		t.Errorf("Expected New, got %s", updated[0].Desc)
	}

	_, err = UpdateTaskDescription(items, 2, "Oops")
	if err == nil {
		t.Errorf("Expected error for out-of-range update, got nil")
	}
}

// Test updating a task status
func TestUpdateTaskStatus(t *testing.T) {
	items := []TodoItem{{Desc: "Task Z", Status: "Not started"}}
	updated, err := UpdateTaskStatus(items, 0, "Completed")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if updated[0].Status != "Completed" {
		t.Errorf("Expected Completed, got %s", updated[0].Status)
	}

	_, err = UpdateTaskStatus(items, 10, "Started")
	if err == nil {
		t.Errorf("Expected error for out-of-range update, got nil")
	}
}

// Test deleting a task
func TestDeleteTask(t *testing.T) {
	items := []TodoItem{
		{Desc: "A", Status: "Not started"},
		{Desc: "B", Status: "Started"},
		{Desc: "C", Status: "Completed"},
	}

	updated, err := DeleteTask(items, 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(updated) != 2 || updated[0].Desc != "A" || updated[1].Desc != "C" {
		t.Errorf("Unexpected result: %+v", updated)
	}

	_, err = DeleteTask(items, 99)
	if err == nil {
		t.Errorf("Expected error for out-of-range delete, got nil")
	}
}
