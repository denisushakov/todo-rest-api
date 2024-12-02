package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/subosito/gotenv"
)

const (
	DefaultPort   = "8080"
	DefaultDBFile = "./storage/scheduler.db"
	DefaultWebDir = "./web"
	MaxTaskLimit  = 50
)

var (
	Port           string
	DBFilePath     string
	WebDirPath     string
	Password       string
	SecretKeyBytes []byte
)

func MustLoad() {

	dir, err := os.Getwd() // current directory
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	if filepath.Base(dir) == "tests" {
		dir = filepath.Dir(dir)
	}

	err = gotenv.Load(absPath(dir, ".env"))
	if err != nil {
		log.Fatalf("env file is not set: %v", err)
	}

	WebDirPath = absPath(dir, DefaultWebDir)

	Port = os.Getenv("TODO_PORT")
	if Port == "" {
		Port = DefaultPort
	}

	DBFilePath = os.Getenv("TODO_DBFILE")
	if DBFilePath == "" {
		DBFilePath = DefaultDBFile
	}
	DBFilePath = absPath(dir, DBFilePath)

	Password = os.Getenv("TODO_PASSWORD")

	secretKey := os.Getenv("TODO_JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("secret key is empty")
	}
	SecretKeyBytes = []byte(secretKey)

}

func absPath(dir, path string) string {
	return filepath.Join(dir, path)
}
