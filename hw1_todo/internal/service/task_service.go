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
		return domain.Task{}, fmt.Errorf("create task: %w", err)
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

	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.Limit == 0 {
		opts.Limit = 20
	}

	filteredTasks, total, err := ts.repo.List(opts)
	if err != nil {
		return domain.ListTaskOutput{}, fmt.Errorf("get list: %w", err)
	}
	if total == 0 {
		return domain.ListTaskOutput{}, fmt.Errorf("get list: %w", domain.ErrorTasksNotFound)
	}

	totalPages := (total + opts.Limit - 1) / opts.Limit
	if opts.Page > totalPages {
		return domain.ListTaskOutput{}, fmt.Errorf("get list: %w", domain.ErrorOptsValidatePage)
	}

	return domain.ListTaskOutput{
		Tasks:      filteredTasks,
		Limit:      opts.Limit,
		Page:       opts.Page,
		TotalPages: totalPages,
	}, nil
}

func (ts *TaskService) UpdateStatus(id domain.TaskID, newStatus domain.TaskStatus) (domain.Task, error) {
	newStatus, err := domain.ParseTaskStatus(string(newStatus))
	if err != nil {
		return domain.Task{}, fmt.Errorf("update status: %w", err)
	}
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
		ID:       task.ID,
		Title:    opts.Title,
		Status:   task.Status,
		Deadline: task.Deadline,
	}
	if !opts.Deadline.IsZero() {
		changedTask.Deadline = opts.Deadline
	}

	changedTask, err = ts.repo.Update(changedTask)
	if err != nil {
		return domain.Task{}, fmt.Errorf("update task: %w", err)
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
	if !deadline.IsZero() && time.Time.Before(deadline, time.Now()) {
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
	if opts.Status != nil {
		if _, err := domain.ParseTaskStatus(string(*opts.Status)); err != nil {
			return fmt.Errorf("validate status: %w", err)
		}
	}
	return nil
}
