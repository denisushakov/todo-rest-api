package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/denisushakov/todo-rest/internal/config"
	"github.com/denisushakov/todo-rest/pkg/models"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrNotFound = errors.New("record not found")
)

type Storage struct {
	db *sql.DB
}

type Search struct {
	Search string
	Date   string
	Active bool
}

func New(storagePath string) (*Storage, error) {

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

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

	return &Storage{db: db}, nil
}

func (s *Storage) SaveTask(task *models.Task) (int64, error) {

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
		return 0, fmt.Errorf("%w", err)
	}

	return id, nil
}

func (s *Storage) GetTasks(search_st *Search) ([]*models.Task, error) {
	var query string
	args := []any{}

	query = "SELECT * FROM scheduler ORDER BY date LIMIT :limit"
	if search_st.Active && search_st.Search != "" {
		query = "SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit"
		args = append(args, sql.Named("search", fmt.Sprintf("%%%s%%", search_st.Search)))
	} else if search_st.Active && search_st.Date != "" {
		query = "SELECT * FROM scheduler WHERE date = :date LIMIT :limit"
		args = append(args, sql.Named("date", search_st.Date))
	}

	args = append(args, sql.Named("limit", config.MaxTaskLimit))

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]*models.Task, 0, 10)

	for rows.Next() {
		var task models.Task
		rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (s *Storage) GetTaskByID(id string) (*models.Task, error) {
	var task models.Task

	query := "SELECT * FROM scheduler WHERE id = ?"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	err = stmt.QueryRow(id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *Storage) UpdateTask(task *models.Task) error {
	query := `UPDATE scheduler SET
		date = :date,
		title = :title,
		comment = :comment,
		repeat = :repeat
	WHERE id = :id`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	res, err := stmt.Exec(
		sql.Named("id", &task.ID),
		sql.Named("date", &task.Date),
		sql.Named("title", &task.Title),
		sql.Named("comment", &task.Comment),
		sql.Named("repeat", &task.Repeat),
	)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if num == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *Storage) DeleteTask(id string) error {
	res, err := s.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if num == 0 {
		return ErrNotFound
	}

	return nil
}
