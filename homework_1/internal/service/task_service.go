package service

import (
	"fmt"
	"time"

	"github.com/vyache31/go-ozon-homeworks/homework_1/internal/domain"
)

type TaskRepository interface {
	Create(task domain.Task) (domain.Task, error)
	List(opts TaskListOptions) ([]domain.Task, int, error)
	Get(id domain.TaskID) (domain.Task, error)
	Update(task domain.Task) (domain.Task, error)
	Delete(id domain.TaskID) error
}

type TaskService struct {
	repo TaskRepository
}

func NewService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (ts *TaskService) Create(input domain.CreateTaskInput) (domain.Task, error) {
	if err := ts.validateTaskFields(input.Title, input.Deadline); err != nil {
		return domain.Task{}, fmt.Errorf("create task %w:", err)
	}

	newTask, err := ts.repo.Create(domain.Task{
		Title:    input.Title,
		Status:   domain.StatusPlanned,
		Deadline: input.Deadline,
	})

	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return newTask, nil
}

func (ts *TaskService) Get(id domain.TaskID) (domain.Task, error) {
	task, err := ts.repo.Get(id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task %d: %w", id, err)
	}

	return task, nil
}

func (ts *TaskService) List(opts TaskListOptions) (domain.ListTaskOutput, error) {
	err := ts.validateListOptions(opts)
	if err != nil {
		return domain.ListTaskOutput{}, fmt.Errorf("get list: %w", err)
	}
	filteredTasks, total, err := ts.repo.List(opts)
	if err != nil {
		return domain.ListTaskOutput{}, fmt.Errorf("get list: %w", err)
	}

	totalPages := total / opts.Limit

	return domain.ListTaskOutput{
		Tasks:      filteredTasks,
		TotalPages: totalPages,
	}, nil
}

func (ts *TaskService) UpdateStatus(id domain.TaskID, newStatus domain.TaskStatus) (domain.Task, error) {
	task, err := ts.repo.Get(id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("update status task %d: %w", id, err)
	}
	if task.Status == newStatus {
		return domain.Task{}, fmt.Errorf("update status task %d: %w", id, domain.ErrorTaskSameStatus)
	}
	if task.Status == domain.StatusDone || task.Status == domain.StatusCanceled {
		return domain.Task{}, fmt.Errorf("update status task %d: %w", id, domain.ErrorTaskLockedStatus)
	}

	updatedTask, err := ts.repo.Update(domain.Task{
		ID:       task.ID,
		Title:    task.Title,
		Status:   newStatus,
		Deadline: task.Deadline,
	})
	if err != nil {
		return domain.Task{}, fmt.Errorf("update status task %d: %w", id, err)
	}

	return updatedTask, nil
}

func (ts *TaskService) Update(id domain.TaskID, opts TaskUpdateOptions) (domain.Task, error) {
	task, err := ts.repo.Get(id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("update task %d: %w", id, err)
	}
	if task.Status == domain.StatusDone || task.Status == domain.StatusCanceled {
		return domain.Task{}, fmt.Errorf("update task %d: %w", id, domain.ErrorTaskLockedStatus)
	}
	if err := ts.validateTaskFields(opts.Title, opts.Deadline); err != nil {
		return domain.Task{}, fmt.Errorf("update task %d: %w", id, err)
	}

	changedTask := domain.Task{
		ID:     task.ID,
		Title:  task.Title,
		Status: task.Status,
	}
	if !opts.Deadline.IsZero() {
		changedTask.Deadline = opts.Deadline
	}
	return changedTask, nil
}

func (ts *TaskService) Delete(id domain.TaskID) error {
	err := ts.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("delete task %d: %w", id, err)
	}

	return nil
}

func (ts *TaskService) validateTaskFields(title string, deadline time.Time) error {
	if title == "" {
		return domain.ErrorTaskValidateTitle
	}
	if time.Time.Before(deadline, time.Now()) || deadline.IsZero() {
		return domain.ErrorTaskValidateDeadline
	}

	return nil
}

func (ts *TaskService) validateListOptions(opts TaskListOptions) error {
	if opts.Page < 0 {
		return fmt.Errorf("%w: page must be more than 0", domain.ErrorOptsValidatePage)
	}
	if opts.Limit < 0 || opts.Limit > 100 {
		return fmt.Errorf("%w: limit must be between 1 and 100", domain.ErrorOptsValidateLimit)
	}

	return nil
}
