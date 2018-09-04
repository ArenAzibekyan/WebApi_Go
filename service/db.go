package service

import (
	"github.com/beevik/guid"
	"github.com/asdine/storm"
	"time"
)

// model
type Task struct {
	Guid      *guid.Guid `json:"guid" storm:"id"`
	Timestamp time.Time  `json:"timestamp"`
	Status    string     `json:"status"`
}

// storm db instance
var (
	taskDb *storm.DB
)

// open database
func OpenTaskDb() error {
	db, err := storm.Open("task.db")
	if err != nil {
		return err
	}

	taskDb = db
	return nil
}

// close database
func CloseTaskDb() error {
	return taskDb.Close()
}

// create new task
func CreateTask(guid *guid.Guid) error {
	taskItem := &Task{
		Guid:      guid,
		Status:    "created",
		Timestamp: time.Now(),
	}

	return taskDb.Save(taskItem)
}

// update task status
func updateStatus(guid *guid.Guid, status string) error {
	taskItem := &Task{
		Guid:      guid,
		Status:    status,
		Timestamp: time.Now(),
	}

	return taskDb.Update(taskItem)
}

func RunTask(guid *guid.Guid) error {
	return updateStatus(guid, "running")
}

func FinishTask(guid *guid.Guid) error {
	return updateStatus(guid, "finished")
}

// get task by guid
func TaskByGuid(guid *guid.Guid) (*Task, error) {
	var taskItem Task

	err := taskDb.One("Guid", guid, &taskItem)
	if err != nil {
		return nil, err
	}

	return &taskItem, nil
}

// all created tasks
func AllTasks() ([]Task, error) {
	var tasks []Task
	err := taskDb.All(&tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// delete task
func DeleteTask(taskItem *Task) error {
	return taskDb.DeleteStruct(taskItem)
}
