package tests

import (
	"net/http/httptest"

	"github.com/denisushakov/todo-rest/internal/config"
	"github.com/denisushakov/todo-rest/pkg/router"

	_ "github.com/mattn/go-sqlite3"
)

func createTestServer() *httptest.Server {
	config.MustLoad()

	router := router.SetupRouter()

	return httptest.NewServer(router)
}
