package main

import (
	"log"
	"net/http"

	"github.com/denisushakov/todo-rest/internal/config"
	"github.com/denisushakov/todo-rest/pkg/router"
)

func main() {
	cfg := config.MustLoad()
	_ = cfg

	port := ":" + config.Port

	router := router.SetupRouter()

	log.Printf("Server is running at %s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	log.Fatalf("server stopped")
}
