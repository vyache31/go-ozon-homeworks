package repository

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/vyache31/go-ozon-homeworks/homework_1/internal/domain"
	"github.com/vyache31/go-ozon-homeworks/homework_1/internal/service"
)

type JSONTaskRepository struct {
	Tasks  map[domain.TaskID]domain.Task
	path   string
	NextID domain.TaskID
}

func NewFileTaskRepository(path string) *JSONTaskRepository {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &JSONTaskRepository{
			Tasks:  make(map[domain.TaskID]domain.Task),
			path:   path,
			NextID: 0,
		}
	} else if err == nil {
		repo := &JSONTaskRepository{}
		if err := uploadJSON(path, repo); err != nil {
			panic(err.Error())
		}
		repo.path = path
		return repo
	} else {
		panic(fmt.Errorf("failed to create a NewFileTaskRepository: %w", err))
	}
}

func (r *JSONTaskRepository) Create(task domain.Task) (domain.Task, error) {
	newTask := domain.Task{
		ID:       r.NextID,
		Title:    task.Title,
		Status:   task.Status,
		Deadline: task.Deadline,
	}
	r.Tasks[newTask.ID] = newTask

	if err := r.saveJSON(); err != nil {
		delete(r.Tasks, r.NextID)

		return domain.Task{}, err
	}
	r.NextID++

	return newTask, nil
}

func (r *JSONTaskRepository) List(opts service.TaskListOptions) ([]domain.Task, int, error) {
	filteredTasks := make([]domain.Task, 0, len(r.Tasks))
	for _, task := range r.Tasks {
		if opts.Search != "" && !strings.Contains(task.Title, opts.Search) {
			continue
		}
		if opts.Overdue && !task.IsOverdue(time.Now()) {
			continue
		}
		if opts.Status != nil && *opts.Status != task.Status {
			continue
		}

		filteredTasks = append(filteredTasks, task)
	}
	total := len(filteredTasks)
	offset := (opts.Page - 1) * opts.Limit

	if offset >= total {
		return []domain.Task{}, total, nil
	}

	end := offset + opts.Limit
	if end > total {
		end = total
	}
	slices.SortFunc(filteredTasks, func(a, b domain.Task) int {
		if result := a.Deadline.Compare(b.Deadline); result != 0 {
			return result
		}

		return cmp.Compare(a.ID, b.ID)
	})

	return filteredTasks[offset:end], total, nil
}

func (r *JSONTaskRepository) Get(id domain.TaskID) (domain.Task, error) {
	task, exists := r.Tasks[id]
	if !exists {
		return domain.Task{}, domain.ErrorTaskNotFound
	}

	return task, nil
}

func (r *JSONTaskRepository) Update(changedTask domain.Task) (domain.Task, error) {
	task, exists := r.Tasks[changedTask.ID]
	if !exists {
		return domain.Task{}, domain.ErrorTaskNotFound
	}
	tmpTask := task
	r.Tasks[changedTask.ID] = changedTask
	if err := r.saveJSON(); err != nil {
		r.Tasks[changedTask.ID] = tmpTask

		return domain.Task{}, err
	}

	return task, nil
}

func (r *JSONTaskRepository) Delete(id domain.TaskID) error {
	task, exists := r.Tasks[id]
	if !exists {
		return domain.ErrorTaskNotFound
	}
	delete(r.Tasks, id)

	if err := r.saveJSON(); err != nil {
		r.Tasks[id] = task
		return err
	}

	return nil
}

func uploadJSON(path string, repo *JSONTaskRepository) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read JSON: %w", err)
	}
	err = json.Unmarshal(file, repo)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

func (r *JSONTaskRepository) saveJSON() error {
	newJSON, err := json.Marshal(r)
	tmpPath := r.path + ".tmp"
	if err != nil {
		return fmt.Errorf("failed to serialize tasks in JSON: %w", err)
	}

	if err := os.WriteFile(tmpPath, newJSON, 0644); err != nil {
		return fmt.Errorf("failed to save temp JSON: %w", err)
	}

	if err := os.Rename(tmpPath, r.path); err != nil {
		return fmt.Errorf("failed to rename temp JSON to main file: %w", err)
	}
	return nil
}
