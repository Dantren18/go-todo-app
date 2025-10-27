package store

import (
	"errors"
	"os"
	"strings"
)

// function to load tasks from a given file path
func LoadTasks(filename string) ([]string, error) {
	//attempt to read from the file
	fileData, err := os.ReadFile(filename)

	//if file does not exist, return an empty task list
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		//return the error if something else went wrong
		return nil, err
	}

	//convert file content to string
	fileContent := string(fileData)

	//split string into slice of tasks using newline as delimiter
	tasks := strings.Split(fileContent, "\n")

	//return the tasks slice and nil error
	return tasks, nil
}

// function to save a slice of tasks back to file
func SaveTasks(filename string, tasks []string) error {
	//join slice of tasks into single string with newlines
	fileContent := strings.Join(tasks, "\n")

	//write the string to the file
	err := os.WriteFile(filename, []byte(fileContent), 0644)

	//return any error that occurs during write
	return err
}

// function to update a task at a given index
func UpdateTask(tasks []string, index int, newValue string) ([]string, error) {
	//check if index is within valid range
	if index < 0 || index >= len(tasks) {
		return tasks, errors.New("update index out of range")
	}

	//replace task at the given index
	tasks[index] = newValue

	//return updated tasks slice
	return tasks, nil
}

// function to delete a task at a given index
func DeleteTask(tasks []string, index int) ([]string, error) {
	//check if index is within valid range
	if index < 0 || index >= len(tasks) {
		return tasks, errors.New("delete index out of range")
	}

	//remove the task by slicing around the index
	tasks = append(tasks[:index], tasks[index+1:]...)

	//return updated slice with task removed
	return tasks, nil
}
