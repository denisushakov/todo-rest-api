package handlers

import "github.com/denisushakov/todo-rest/pkg/models"

type TaskSaver interface {
	SaveTask(task *models.Task) (int64, error)
}

type TaskGetter interface {
	GetTasks(search string) ([]*models.Task, error)
	GetTaskByID(id string) (*models.Task, error)
}

type TaskUpdater interface {
	UpdateTask(task *models.Task) error
}

type TaskConditionUpdater interface {
	MarkTaskCompleted(id string) error
}

type TaskRemover interface {
	DeleteTask(id string) error
}
