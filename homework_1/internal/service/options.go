package service

import (
	"time"

	"github.com/vyache31/go-ozon-homeworks/homework_1/internal/domain"
)

type TaskListOptions struct {
	Status  *domain.TaskStatus // если status не задан, то значение будет nil
	Page    int
	Limit   int
	Overdue bool
	Search  string
}

type TaskUpdateOptions struct {
	Title    string
	Status   domain.TaskStatus
	Deadline time.Time
}
