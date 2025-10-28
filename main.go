package main

import (
	"GoCourse/store"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//generate a unique trace ID using current timestamp
	traceID := fmt.Sprintf("%d", time.Now())
	ctx := context.WithValue(context.Background(), "traceID", traceID)

	//declaring pointers/variables/flags
	var taskPtr *string
	taskPtr = flag.String("task", "", "The task you want to add or use for updating")

	var updateIndexPtr *int
	updateIndexPtr = flag.Int("update", 0, "The task number to update")

	var deleteIndexPtr *int
	deleteIndexPtr = flag.Int("delete", 0, "The task number to delete")

	flag.Parse()

	fmt.Println("Update index:", *updateIndexPtr)
	fmt.Println("Delete index:", *deleteIndexPtr)

	tasks, err := store.LoadTasks("tasks.txt")

	if err != nil {
		slog.ErrorContext(ctx, "LoadTasks failed", "err", err)
		return
	}

	//update task if -update flag was used
	if *updateIndexPtr > 0 {
		if *taskPtr == "" {
			fmt.Println("Please provide a task with -task to update an item.")
			return
		}
		tasks, err = store.UpdateTask(tasks, *updateIndexPtr-1, *taskPtr)
		if err != nil {
			slog.ErrorContext(ctx, "UpdateTask failed", "err", err)
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Updated task #%d\n", *updateIndexPtr)

		//delete task if -delete flag was used
	} else if *deleteIndexPtr > 0 {
		tasks, err = store.DeleteTask(tasks, *deleteIndexPtr-1)
		if err != nil {
			slog.ErrorContext(ctx, "DeleteTask failed", "err", err)
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Deleted task #%d\n", *deleteIndexPtr)

		//otherwise just add a new task
	} else {
		if *taskPtr == "" {
			fmt.Println("Please provide a task using -task flag.")
			return
		}
		tasks = append(tasks, *taskPtr)
		fmt.Println("Added new task.")
	}

	//print the full contents of the slice
	fmt.Println("Your To-Do List:")
	for i := 0; i < len(tasks); i++ {
		fmt.Printf("%d. %s\n", i+1, tasks[i])
	}

	//save updated tasks back to file using store package
	err = store.SaveTasks("tasks.txt", tasks)

	//error handling if save fails
	if err != nil {
		slog.ErrorContext(ctx, "SaveTasks failed", "err", err)
	} else {
		slog.InfoContext(ctx, "Tasks saved to disk", "count", len(tasks))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Press Ctrl+C to stop running the todo list")

	<-stop

	slog.InfoContext(ctx, "Received shutdown signal, exiting...")
}
