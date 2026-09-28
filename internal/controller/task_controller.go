package controller

import (
	"github.com/vyache31/todo/internal/domain"
	"github.com/vyache31/todo/internal/service"
)

type TaskService interface {
	CreateTask(input domain.CreateTaskInput) error
	List(opts service.TaskListOptions) ([]domain.Task, error)
	Get(id domain.TaskID) (domain.Task, error)
	UpdateStatus(id domain.TaskID, status domain.TaskStatus) error
	Update(opts service.TaskChangeOptions) error
	Delete(id domain.TaskID) error
}

type TaskController struct {
	service TaskService
}
