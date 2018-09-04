package service

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/beevik/guid"
	"github.com/asdine/storm"
	"encoding/json"
	"time"
)

// POST /task
func CreateMethod(w http.ResponseWriter, r *http.Request) {
	taskGuid := guid.New()
	err := CreateTask(taskGuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	go runAndFinish(taskGuid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskGuid.String())
}

func runAndFinish(guid *guid.Guid) {
	RunTask(guid)
	time.Sleep(time.Minute * 2)
	FinishTask(guid)
}

// GET /task/{id}
func GetByIdMethod(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	taskGuid, err := guid.ParseString(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	taskItem, err := TaskByGuid(taskGuid)
	if err != nil {
		if err == storm.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Timestamp time.Time `json:"timestamp"`
		Status    string    `json:"status"`
	}{
		taskItem.Timestamp, taskItem.Status,
	})
}

// GET /task
func GetAllMethod(w http.ResponseWriter, r *http.Request) {
	tasks, err := AllTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ta := make([]struct {
		Guid      string    `json:"guid"`
		Timestamp time.Time `json:"timestamp"`
		Status    string    `json:"status"`
	}, len(tasks))
	for i, t := range tasks {
		ta[i].Guid = t.Guid.String()
		ta[i].Timestamp = t.Timestamp
		ta[i].Status = t.Status
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ta)
}

// DELETE /task/{id}
func DeleteMethod(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	taskGuid, err := guid.ParseString(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	taskItem, err := TaskByGuid(taskGuid)
	if err != nil {
		if err == storm.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = DeleteTask(taskItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
