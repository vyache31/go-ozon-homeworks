package domain

import (
	"time"
)

type TaskID uint

type TaskStatus string

const (
	StatusPlanned    TaskStatus = "planned"
	StatusInProgress TaskStatus = "in_progress"
	StatusCanceled   TaskStatus = "canceled"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID       TaskID
	Title    string
	Status   TaskStatus
	Deadline time.Time
}

func (t Task) IsOverdue(tm time.Time) bool {
	if t.Deadline.Before(tm) &&
		t.Status != StatusDone &&
		t.Status != StatusCanceled {
		return true
	}
	return false
}

// при расширении параметров, передаваемых при создании задач,
// нам нужно лишь добавить их в DTO,
// не меняя аргументы у метода TaskService interface
type CreateTaskInput struct {
	Title    string
	Deadline time.Time
}
