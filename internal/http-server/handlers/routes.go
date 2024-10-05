package handlers

import (
	mwAuth "github.com/denisushakov/todo-rest/internal/http-server/middleware/auth"
	"github.com/denisushakov/todo-rest/internal/scheduler"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router *chi.Mux, scheduler *scheduler.Scheduler) {
	router.Get("/api/nextdate", GetNextDate)
	router.Post("/api/signin", LoginHandler)

	router.Route("/api", func(r chi.Router) {
		r.Use(mwAuth.Auth)

		r.Get("/tasks", GetTasks(scheduler))
		r.Post("/task", SaveTask(scheduler))

		r.Get("/task", GetTaskByID(scheduler))

		r.Put("/task", UpdateTask(scheduler))

		r.Post("/task/done", MarkTaskCompleted(scheduler))

		r.Delete("/task", DeleteTask(scheduler))
	})
}
