package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/denisushakov/todo-rest.git/internal/scheduler"
)

func GetNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nowDate, err := time.Parse("20060102", now)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("time cannot pasre: %s", err)
	}

	newDate, err := scheduler.NextDate(nowDate, date, repeat)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("new date not created: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(newDate))
}
