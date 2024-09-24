package handlers

import "github.com/denisushakov/todo-rest.git/pkg/models"

type TaskSaver interface {
	SaveTask(task *models.Task) (int64, error)
}

type TaskGetter interface {
	GetTasks(search string) ([]*models.Task, error)
}
