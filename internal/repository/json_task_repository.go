package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vyache31/todo/internal/domain"
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
		panic(err.Error())
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

func (r *JSONTaskRepository) List() ([]domain.Task, error) {
	sliceTasks := make([]domain.Task, 0, len(r.Tasks))
	for _, task := range r.Tasks {
		sliceTasks = append(sliceTasks, task)
	}

	return sliceTasks, nil
}

func (r *JSONTaskRepository) Get(id domain.TaskID) (domain.Task, error) {
	task, exists := r.Tasks[id]
	if !exists {
		return domain.Task{}, fmt.Errorf("task with id: %d does not exists", id)
	}

	return task, nil
}

func (r *JSONTaskRepository) UpdateStatus(id domain.TaskID, status domain.TaskStatus) (domain.Task, error) {
	task, exists := r.Tasks[id]
	if !exists {
		return domain.Task{}, fmt.Errorf("task with id: %d does not exists", id)
	}
	tmpStatus := task.Status
	task.Status = status
	r.Tasks[id] = task
	if err := r.saveJSON(); err != nil {
		task.Status = tmpStatus
		r.Tasks[id] = task

		return domain.Task{}, err
	}

	return task, nil
}

func (r *JSONTaskRepository) Delete(id domain.TaskID) error {
	task, exists := r.Tasks[id]
	if !exists {
		return fmt.Errorf("task with id: %d does not exists", id)
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
		return err
	}
	err = json.Unmarshal(file, repo)
	if err != nil {
		return err
	}

	return nil
}

func (r *JSONTaskRepository) saveJSON() error {
	newJSON, err := json.Marshal(r)
	tmpPath := r.path + ".tmp"
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmpPath, newJSON, 0644); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, r.path); err != nil {
		return err
	}
	return nil
}
