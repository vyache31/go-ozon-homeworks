package main

import (
	"github.com/vyache31/go-ozon-homeworks/homework_1/internal/controller"
	"github.com/vyache31/go-ozon-homeworks/homework_1/internal/repository"
	"github.com/vyache31/go-ozon-homeworks/homework_1/internal/service"
)

func main() {
	repo := repository.NewFileTaskRepository("tasks.json")
	ts := service.NewService(repo)
	tc := controller.NewTaskController(ts)
	tc.Run()
}
