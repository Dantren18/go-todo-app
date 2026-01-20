package store

import (
	"fmt"
	"log"
)

type TaskCommandType int

// list of all the possible types of commands to send to actor
const (
	GetAllTasks TaskCommandType = iota
	CreateTask
	UpdateTask
	DeleteaTask
)

type TaskCommand struct {
	Type       TaskCommandType // what kind of command
	NewTask    TodoItem
	Index      int
	Desc       string
	Status     string
	ResponseCh chan []TodoItem
}

// commandCh is a shared channel where commands are sent to the actor
var commandCh chan TaskCommand

func StartStoreActor() {
	commandCh = make(chan TaskCommand)

	// launch a goroutine so the actor runs in the background
	go func() {
		// load the tasks from the file into memory when app starts
		tasks, err := LoadTasks("tasks.json")
		if err != nil {
			log.Fatal("Failed to load tasks:", err)
		}

		// listen for new commands sent on the channel
		for cmd := range commandCh {
			fmt.Println("Processing command:", cmd.Type)

			switch cmd.Type {
			case GetAllTasks:
				cmd.ResponseCh <- tasks

			case CreateTask:
				// validate status for new tasks - if invalid, respond with current list and do not save
				if cmd.NewTask.Status != "Not Started" && cmd.NewTask.Status != "Started" && cmd.NewTask.Status != "Completed" {
					fmt.Println("Failed to create task: invalid status", cmd.NewTask.Status)
					cmd.ResponseCh <- tasks
				} else {
					tasks = append(tasks, cmd.NewTask)
					SaveTasks("tasks.json", tasks)
					fmt.Println("Created new task, lists of tasks:", tasks)
					cmd.ResponseCh <- tasks
				}

			case UpdateTask:
				if cmd.Index >= 0 && cmd.Index < len(tasks) {
					var err error
					if cmd.Desc != "" {
						updated, e := UpdateTaskDescription(tasks, cmd.Index, cmd.Desc)
						if e != nil {
							fmt.Println("Failed to update description:", e)
						} else {
							tasks = updated
							err = e
						}
					}
					if cmd.Status != "" {
						updated, e := UpdateTaskStatus(tasks, cmd.Index, cmd.Status)
						if e != nil {
							fmt.Println("Failed to update status:", e)
						} else {
							tasks = updated
							err = e
						}
					}
					if err == nil {
						SaveTasks("tasks.json", tasks)
					}
				}
				cmd.ResponseCh <- tasks

			case DeleteaTask:
				if cmd.Index >= 0 && cmd.Index < len(tasks) {
					updated, e := DeleteTask(tasks, cmd.Index)
					if e != nil {
						fmt.Println("Failed to delete task:", e)
					} else {
						tasks = updated
						SaveTasks("tasks.json", tasks)
					}
				}
				cmd.ResponseCh <- tasks
			}
		}

	}()
}

// this function is used to send a request to actor
func SendCommand(cmd TaskCommand) []TodoItem {
	// make a new channel where the actor will send the result  send the command to the actor
	cmd.ResponseCh = make(chan []TodoItem)
	commandCh <- cmd

	return <-cmd.ResponseCh
}
