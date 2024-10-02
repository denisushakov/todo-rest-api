package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/denisushakov/todo-rest/internal/storage/sqlite"
	"github.com/denisushakov/todo-rest/pkg/models"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	response := ErrorResponse{
		Error: err.Error(),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("error marshal JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}

func SaveTask(taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		id, err := taskSaver.SaveTask(&task)
		if err != nil {
			writeErrorResponse(w, err, http.StatusBadRequest)
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
			writeErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
	}
}

func GetTaskByID(taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		if id == "" {
			http.Error(w, `{"error": "id not specified"}`, http.StatusBadRequest)
			return
		}

		task, err := taskGetter.GetTaskByID(id)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				writeErrorResponse(w, err, http.StatusNotFound)
			default:
				writeErrorResponse(w, err, http.StatusInternalServerError)
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
			writeErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		if err := taskUpdater.UpdateTask(&task); err != nil {
			switch {
			case errors.Is(err, sqlite.ErrNotFound):
				writeErrorResponse(w, err, http.StatusNotFound)
			default:
				writeErrorResponse(w, err, http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

func MarkTaskCompleted(taskConditionUpdater TaskConditionUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		if id == "" {
			http.Error(w, `{"error": "id not specified"}`, http.StatusBadRequest)
			return
		}

		if err := taskConditionUpdater.MarkTaskCompleted(id); err != nil {
			writeErrorResponse(w, err, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

func DeleteTask(taskRemover TaskRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		if id == "" {
			http.Error(w, `{"error": "id not specified"}`, http.StatusBadRequest)
			return
		}

		if err := taskRemover.DeleteTask(id); err != nil {
			writeErrorResponse(w, err, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
