package main

import (
	"log"
	"net/http"
	"os"

	"github.com/denisushakov/todo-rest.git/internal/config"
	"github.com/denisushakov/todo-rest.git/internal/http-server/handlers"
	"github.com/denisushakov/todo-rest.git/internal/scheduler"
	"github.com/denisushakov/todo-rest.git/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	//fmt.Println(cfg)
	_ = cfg

	webDir := config.WebDir
	port := ":" + config.GetPort()

	storage, err := sqlite.New(config.DBFilePath)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		os.Exit(1)
	}
	scheduler := scheduler.NewScheduler(storage)

	router := chi.NewRouter()

	router.Use(middleware.URLFormat)

	router.Handle("/*", http.FileServer(http.Dir(webDir)))

	router.Get("/api/nextdate", handlers.GetNextDate)

	router.Route("/api", func(r chi.Router) {
		r.Get("/tasks", handlers.GetTasks(scheduler))
		r.Post("/task", handlers.SaveTask(scheduler))

		r.Get("/task", handlers.GetTask(scheduler))

		r.Put("/task", handlers.UpdateTask(scheduler))
	})

	log.Printf("Server is running at %s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	log.Fatalf("server stopped")
}
