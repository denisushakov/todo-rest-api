package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/denisushakov/todo-rest.git/internal/storage/sqlite"
	"github.com/denisushakov/todo-rest.git/pkg/models"
)

func SaveTask(taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.Task

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"failed to decode JSON"}`, http.StatusBadRequest)
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

func GetTask(taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		if id == "" {
			http.Error(w, `{"error": "id not specified"}`, http.StatusBadRequest)
			return
		}

		task, err := taskGetter.GetTask(id)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				http.Error(w, `{"error": "task not found"}`, http.StatusNotFound)
			default:
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				//http.Error(w, fmt.Sprintf(`{"error":%s"}`, err.Error()), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(task)
	}
}

func UpdateTask(taskUpdater TaskUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, `{"error":"failed to decode JSON"}`, http.StatusBadRequest)
			return
		}

		if err := taskUpdater.UpdateTask(&task); err != nil {
			switch {
			case errors.Is(err, sqlite.ErrNotFound):
				http.Error(w, `{"error":"record not found"}`, http.StatusNotFound)
			default:
				//http.Error(w, `{"error":"bad"}`, http.StatusInternalServerError)
				http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
				//http.Error(w, fmt.Sprintf(`{"error":%s"}`, err.Error()), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
