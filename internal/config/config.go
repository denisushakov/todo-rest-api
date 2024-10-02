package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/subosito/gotenv"
)

const (
	DefaultPort  = "8080"
	DBFile       = "./storage/scheduler.db"
	WebDir       = "./web"
	MaxTaskLimit = 50
)

var (
	Port           string
	DBFilePath     string
	Password       string
	SecretKeyBytes []byte
)

type Config struct {
}

func MustLoad() *Config {

	err := gotenv.Load()
	if err != nil {
		log.Fatalf("env file is not set: %v", err)
	}
	Port = os.Getenv("TODO_PORT")
	DBFilePath = os.Getenv("TODO_DBFILE")
	Password = os.Getenv("TODO_PASSWORD")

	secretKey := os.Getenv("TODO_JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("secret key is empty")
	}
	SecretKeyBytes = []byte(secretKey)

	var cfg Config

	return &cfg
}

func GetPort() string {
	port := Port
	if port == "" {
		port = DefaultPort
	}
	return port
}

func GetDBFilePath(defaultDBName string) string {
	dbFilePath := DBFilePath
	if dbFilePath != "" {
		return dbFilePath
	}
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}

	parentDir := filepath.Dir(filepath.Dir(executablePath))

	return filepath.Join(parentDir, defaultDBName)
}
