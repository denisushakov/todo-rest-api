package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/denisushakov/todo-rest.git/pkg/scheduler"
)

type Request struct {
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type TaskSaver interface {
	SaveTask(task *scheduler.Task) (int64, error)
}

func SaveTask(taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req scheduler.Task

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"Failed to decode JSON"}`, http.StatusBadRequest)
			return
		}

		id, err := taskSaver.SaveTask(&req)
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}
