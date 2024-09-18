package main

import (
	"log"
	"net/http"
	"os"

	"github.com/denisushakov/todo-rest.git/internal/config"
	"github.com/denisushakov/todo-rest.git/internal/http-server/handlers"
	"github.com/denisushakov/todo-rest.git/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()

	//fmt.Println(cfg)
	_ = cfg

	webDir := config.WebDir
	port := ":" + config.GetPort()

	/*db, err := storage.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()*/
	storage, err := sqlite.New(config.DBFilePath)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		os.Exit(1)
	}

	//_ = storage

	router := chi.NewRouter()

	router.Handle("/", http.FileServer(http.Dir(webDir)))

	router.Get("/api/nextdate", handlers.GetNextDate)

	router.Post("/api/task", handlers.SaveTask(storage))

	log.Printf("Server is running at %s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	log.Fatalf("server stopped")
}
