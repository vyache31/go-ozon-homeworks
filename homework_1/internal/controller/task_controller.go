package controller

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

func NewTaskController(service TaskService) *TaskController {
	return &TaskController{
		service: service,
	}
}

func (tc *TaskController) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Таск-трекер запущен")
	for scanner.Scan() {
		tc.handleCommand(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error: %s", err)
	}
}

func (tc *TaskController) handleCommand(text string) {
	commandArgs := strings.Fields(text)
	switch commandArgs[0] {
	case "help":
		tc.handleHelp()
	case "add":
		tc.handleAdd(commandArgs)
	}
}

func (tc *TaskController) handleHelp() {
	cmdsDescription := `Список доступных команд:
	
1. Добавить задачу с указанием заголовка и дедлайна
	add <title> <deadline/2006-01-02>
	
2. Получить список задач с опциональными параметрами
	list <parameter value>...
		status 		- фильтр по статусу (eg. planned; in_progress; canceled; done; overdue)
		page и limit 	- постраничная пагинация (по умолчанию page = 1, limit = 20)
		search 		- поиск задач по подстроке в заголовке (без учёта регистра)

3. Получить конкретную задачу
	task <id>

4. Изменить статус задачи
	chgstatus <id> <new_status>

5. Отредактировать задачу
	chgtask <id> <title> и/или <deadline>

6. Удалить задачу
	deltask

7. Список команд
	help

8. Завершить работу
	exit
`
	fmt.Println(cmdsDescription)
}

func (tc *TaskController) handleAdd(cmdArgs []string) {
	//
}
