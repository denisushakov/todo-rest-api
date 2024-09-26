package handlers

import "github.com/denisushakov/todo-rest.git/pkg/models"

type TaskSaver interface {
	SaveTask(task *models.Task) (int64, error)
}

type TaskGetter interface {
	GetTasks(search string) ([]*models.Task, error)
	GetTask(id string) (*models.Task, error)
}

type TaskUpdater interface {
	UpdateTask(task *models.Task) error
}
