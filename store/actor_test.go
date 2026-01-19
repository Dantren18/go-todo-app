package store

import (
	"os"
	"strconv"
	"sync"
	"testing"
)

func TestConcurrentTaskCreation(t *testing.T) {
	StartStoreActor()

	var wg sync.WaitGroup // create a wait group to manage goroutines
	total := 100

	for i := 0; i < total; i++ {
		wg.Add(1) // tell the wait group we're starting a new goroutine

		go func(i int) {
			defer wg.Done() // signal that this goroutine is finished when it ends

			SendCommand(TaskCommand{
				Type: CreateTask,
				NewTask: TodoItem{
					Description: "Task #" + strconv.Itoa(i),
					Status:      "Not Started",
				},
			})
		}(i)
	}

	wg.Wait()

	tasks := SendCommand(TaskCommand{
		Type: GetAllTasks, // ask the actor to return all tasks
	})

	if len(tasks) != total {
		t.Errorf("expected %d tasks, got %d", total, len(tasks)) // check if all tasks were created
	}

	defer os.Remove("tasks.json")
}
