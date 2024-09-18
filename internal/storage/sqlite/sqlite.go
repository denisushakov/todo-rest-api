package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/denisushakov/todo-rest.git/internal/config"
	"github.com/denisushakov/todo-rest.git/pkg/scheduler"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {

	if storagePath == "" {
		appPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(appPath)
		storagePath = filepath.Join(filepath.Dir(appPath), config.DBFile)
	}

	//_, err := os.Stat(storagePath)
	//install := os.IsNotExist(err)

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	//if install {
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS scheduler(
			id 		INTEGER PRIMARY KEY AUTOINCREMENT,
			date 	TEXT NOT NULL,
			title 	TEXT NOT NULL,
			comment TEXT,
			repeat 	TEXT NOT NULL
				CHECK (length(repeat) <= 128));
	CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date)
	`)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
	}

	if _, err := stmt.Exec(); err != nil {
		return nil, fmt.Errorf("database not opened: %w", err)
	}
	//}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveTask(task *scheduler.Task) (int64, error) {

	stmt, err := s.db.Prepare("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}
