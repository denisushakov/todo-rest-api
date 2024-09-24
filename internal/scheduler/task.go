package scheduler

import (
	"fmt"
	"time"

	"github.com/denisushakov/todo-rest.git/internal/storage/sqlite"
	"github.com/denisushakov/todo-rest.git/pkg/models"
)

/*
	type Task struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}
*/
type Scheduler struct {
	Storage *sqlite.Storage
}

func NewScheduler(dataBase *sqlite.Storage) *Scheduler {
	return &Scheduler{
		Storage: dataBase,
	}
}

type TaskScheduler interface {
	SaveTask(task *models.Task) (int64, error)
}

func (s *Scheduler) SaveTask(task *models.Task) (int64, error) {
	var now = time.Now().Truncate(24 * time.Hour)
	var nextdate string

	if task.Date == "" {
		nextdate = now.Format("20060102")
	} else {
		date, err := time.Parse("20060102", task.Date)
		if err != nil {
			return 0, fmt.Errorf("invalid date format")
		}
		nextdate = date.Format("20060102")
		if date.Before(now) {
			if task.Repeat == "" {
				nextdate = now.Format("20060102")
			} else {
				nextdate, err = NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return 0, err
				}
			}
		}
	}
	task.Date = nextdate

	id, err := s.Storage.SaveTask(task)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Scheduler) GetTasks(search string) ([]*models.Task, error) {

	var sr_st sqlite.Search
	if search != "" {
		sr_st.Active = true
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			sr_st.Search = search
		} else {
			sr_st.Date = date.Format("20060102")
		}
	}

	tasks, err := s.Storage.GetTasks(&sr_st)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
