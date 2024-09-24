package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/denisushakov/todo-rest.git/pkg/models"
)

type Request struct {
	ID      string `json:"id"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func SaveTask(taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.Task

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"Failed to decode JSON"}`, http.StatusBadRequest)
			return
		}

		if req.Title == "" {
			http.Error(w, `{"error":"Empty title field"}`, http.StatusBadRequest)
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

func GetTasks(taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")

		tasks, err := taskGetter.GetTasks(search)
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
	}
}
