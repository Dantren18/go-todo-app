package store

import (
	"errors"
	"os"
	"strings"
)

// function to load tasks from a given file path
func LoadTasks(filename string) ([]string, error) {
	fileData, err := os.ReadFile(filename)

	//if file does not exist, return an empty task list
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	fileContent := string(fileData)
	tasks := strings.Split(fileContent, "\n")

	return tasks, nil
}

func SaveTasks(filename string, tasks []string) error {
	//join slice of tasks into single string with newlines
	fileContent := strings.Join(tasks, "\n")

	//write the string to the file
	err := os.WriteFile(filename, []byte(fileContent), 0644)
	return err
}

func UpdateTask(tasks []string, index int, newValue string) ([]string, error) {
	if index < 0 || index >= len(tasks) {
		return tasks, errors.New("update index out of range")
	}

	tasks[index] = newValue
	return tasks, nil
}

// function to delete a task at a given index
func DeleteTask(tasks []string, index int) ([]string, error) {
	if index < 0 || index >= len(tasks) {
		return tasks, errors.New("delete index out of range")
	}

	//remove the task by slicing around the index
	tasks = append(tasks[:index], tasks[index+1:]...)
	return tasks, nil
}
