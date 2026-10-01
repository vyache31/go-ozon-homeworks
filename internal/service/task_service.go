package service

import "github.com/vyache31/todo/internal/domain"

type TaskRepository interface {
	Create(task domain.Task) (domain.Task, error)
	List() ([]domain.Task, error)
	Get(id domain.TaskID) (domain.Task, error)
	Update(task domain.Task) (domain.Task, error)
	Delete(id domain.TaskID) error
}

type Service struct {
	repo TaskRepository
}
