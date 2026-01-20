package store

import (
	"os"
	"strconv"
	"sync"
	"testing"
)

func TestConcurrentTaskCreation(t *testing.T) {
	// Remove tasks.json before starting the test to avoid leftover tasks
	_ = os.Remove("tasks.json")
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
		Type: GetAllTasks,
	})

	if len(tasks) != total {
		t.Errorf("expected %d tasks, got %d", total, len(tasks)) // check if all tasks were created
	}

	defer os.Remove("tasks.json")
}

func TestConcurrentTaskCreationParallel(t *testing.T) {
	StartStoreActor()

	total := 50 // number of parallel creators; keep this modest for a beginner machine

	// Create subtests. Each subtest is marked with t.Parallel so the Go test
	// runner runs them concurrently.
	for i := 0; i < total; i++ {
		taskNumber := i
		//create a subtest for each task creation, with a unique name
		t.Run("create-"+strconv.Itoa(taskNumber), func(t *testing.T) {
			t.Parallel()

			SendCommand(TaskCommand{
				Type: CreateTask,
				NewTask: TodoItem{
					Description: "Parallel task " + strconv.Itoa(taskNumber),
					Status:      "Not Started",
				},
			})
		})
	}

	// After all subtests are finished the test will continue to run here.
	tasks := SendCommand(TaskCommand{Type: GetAllTasks})
	if len(tasks) != total {
		t.Errorf("expected %d tasks, got %d", total, len(tasks))
	}

	defer os.Remove("tasks.json")
}

func TestConcurrentTaskUpdateParallel(t *testing.T) {
	_ = os.Remove("tasks.json")
	StartStoreActor()

	// First, create some tasks to update
	total := 20
	for i := 0; i < total; i++ {
		SendCommand(TaskCommand{
			Type: CreateTask,
			NewTask: TodoItem{
				Description: "Task to update " + strconv.Itoa(i),
				Status:      "Not Started",
			},
		})
	}

	// Now, update each task in parallel subtests
	for i := 0; i < total; i++ {
		taskNum := i
		t.Run("update-"+strconv.Itoa(taskNum), func(t *testing.T) {
			t.Parallel()
			SendCommand(TaskCommand{
				Type:   UpdateTask,
				Index:  taskNum,
				Status: "Completed",
			})
		})
	}

	tasks := SendCommand(TaskCommand{Type: GetAllTasks})
	for i, task := range tasks {
		if task.Status != "Completed" {
			t.Errorf("task %d not updated, got status %s", i, task.Status)
		}
	}
	defer os.Remove("tasks.json")
}

// func TestConcurrentTaskDeleteParallel(t *testing.T) {
// 	_ = os.Remove("tasks.json")
// 	StartStoreActor()

// 	// First, create some tasks to delete
// 	total := 20
// 	for i := 0; i < total; i++ {
// 		SendCommand(TaskCommand{
// 			Type: CreateTask,
// 			NewTask: TodoItem{
// 				Description: "Task to delete " + strconv.Itoa(i),
// 				Status:      "Not Started",
// 			},
// 		})
// 	}

// 	// Delete tasks in parallel (from last to first to avoid index shifting)
// 	for i := total - 1; i >= 0; i-- {
// 		taskNum := i
// 		t.Run("delete-"+strconv.Itoa(taskNum), func(t *testing.T) {
// 			t.Parallel()
// 			SendCommand(TaskCommand{
// 				Type:  DeleteaTask,
// 				Index: taskNum,
// 			})
// 		})
// 	}

// 	// Check that all tasks were deleted
// 	tasks := SendCommand(TaskCommand{Type: GetAllTasks})
// 	if len(tasks) != 0 {
// 		t.Errorf("expected 0 tasks after delete, got %d", len(tasks))
// 	}
// 	defer os.Remove("tasks.json")
// }

func TestConcurrentTaskGetParallel(t *testing.T) {
	_ = os.Remove("tasks.json")
	StartStoreActor()

	// Create some tasks to get
	total := 10
	for i := 0; i < total; i++ {
		SendCommand(TaskCommand{
			Type: CreateTask,
			NewTask: TodoItem{
				Description: "Task number " + strconv.Itoa(i),
				Status:      "Not Started",
			},
		})
	}

	// Get each task in parallel subtests
	for i := 0; i < total; i++ {
		taskNum := i
		t.Run("get-"+strconv.Itoa(taskNum), func(t *testing.T) {
			t.Parallel()
			tasks := SendCommand(TaskCommand{Type: GetAllTasks})
			if len(tasks) != total {
				t.Errorf("expected %d tasks, got %d", total, len(tasks))
			}
		})
	}
	defer os.Remove("tasks.json")
}
