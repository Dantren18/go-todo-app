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
				tasks = append(tasks, cmd.NewTask)
				SaveTasks("tasks.json", tasks)
				fmt.Println("Create task, lists of tasks:", tasks)

				cmd.ResponseCh <- tasks

			case UpdateTask:
				if cmd.Index >= 0 && cmd.Index < len(tasks) {
					if cmd.Desc != "" {
						tasks[cmd.Index].Description = cmd.Desc
					}
					if cmd.Status != "" {
						tasks[cmd.Index].Status = cmd.Status
					}
					SaveTasks("tasks.json", tasks)
				}
				cmd.ResponseCh <- tasks

			case DeleteaTask:
				if cmd.Index >= 0 && cmd.Index < len(tasks) {
					tasks = append(tasks[:cmd.Index], tasks[cmd.Index+1:]...)
					SaveTasks("tasks.json", tasks)
				}
				cmd.ResponseCh <- tasks
			}
		}

	}()
}

// this function is used to send a request to actor
func SendCommand(cmd TaskCommand) []TodoItem {
	// make a new channel where the actor will send the result
	cmd.ResponseCh = make(chan []TodoItem)
	// send the command to the actor
	commandCh <- cmd

	return <-cmd.ResponseCh
}
