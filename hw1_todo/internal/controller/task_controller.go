package controller

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vyache31/go-ozon-homeworks/homework_1/internal/domain"
	"github.com/vyache31/go-ozon-homeworks/homework_1/internal/service"
)

type TaskService interface {
	Create(input domain.CreateTaskInput) (domain.Task, error)
	List(opts service.TaskListOptions) (domain.ListTaskOutput, error)
	Get(id domain.TaskID) (domain.Task, error)
	UpdateStatus(id domain.TaskID, status domain.TaskStatus) (domain.Task, error)
	Update(id domain.TaskID, opts service.TaskUpdateOptions) (domain.Task, error)
	Delete(id domain.TaskID) error
}

type CommandParam struct {
	Cmd           string
	RequiredValue bool
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
	fmt.Println("")
	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		tc.handleCommand(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error: %s", err)
	}
}

func (tc *TaskController) handleCommand(text string) {
	commandArgs := strings.Fields(text)
	if len(commandArgs) < 1 {
		fmt.Println("Введите команду")
		return
	}
	fmt.Println("")
	switch commandArgs[0] {
	case "help":
		tc.handleHelp()
	case "add":
		tc.handleAdd(commandArgs[1:])
	case "task":
		tc.handleGet(commandArgs[1:])
	case "list":
		tc.handleList(commandArgs[1:])
	case "chstatus":
		tc.handleChStatus(commandArgs[1:])
	case "chtask":
		tc.handleChTask(commandArgs[1:])
	case "delTask":
		tc.handleDelTask(commandArgs[1:])
	case "exit":
		tc.handleExit()
	}
	fmt.Println("")
}

func (tc *TaskController) handleHelp() {
	cmdsDescription := `Список доступных команд:
	
1. Добавить задачу с указанием заголовка и дедлайна
	add <title> <deadline/2006-01-02>
	
2. Получить список задач с опциональными параметрами
	list <parameter value>...
		status 		- фильтр по статусу (eg. planned; in_progress; canceled; done)
		page и limit 	- постраничная пагинация (по умолчанию page = 1, limit = 20)
		search 		- поиск задач, для многословного значения используйте " " (без учёта регистра)
		overdue 	- параметр без значения, для вывода только просроченных задач

3. Получить конкретную задачу
	task <id>

4. Изменить статус задачи
	chstatus parameter value>...
		title 	- новое название, для многословного значения используйте " " (обязательный параметр)
		dl 	- новый дедлайн (формат 2006-01-02)

5. Отредактировать задачу
	chtask <id> <title> и/или <deadline>

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
	if len(cmdArgs) < 2 {
		fmt.Println("Ошибка: недостаточно аргументов")
		return
	}
	availableParams := []CommandParam{{"title", true}, {"dl", true}}
	cmdParams, err := tc.parseParams(strings.Join(cmdArgs, ""), availableParams)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	deadline, err := time.Parse(time.DateOnly, cmdParams["dl"])
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		fmt.Printf("Проверьте, что дата соответствует формату - %s\n", time.DateOnly)
		return
	}
	createInput := domain.CreateTaskInput{Title: cmdArgs[0], Deadline: deadline}
	task, err := tc.service.Create(createInput)
	if err != nil {
		fmt.Printf("Ошибка: не удалось создать задачу: %v\n", err)
		return
	}

	fmt.Printf("Задача №%d создана!\n", task.ID)
}

func (tc *TaskController) handleGet(cmdArgs []string) {
	if len(cmdArgs) < 1 {
		fmt.Println("Ошибка: недостаточно аргументов")
		return
	}
	value, err := strconv.ParseUint(cmdArgs[0], 10, 64)
	if err != nil {
		fmt.Printf("Ошибка: невалидный ID: %v\n", err)
		return
	}

	taskID := domain.TaskID(value)
	task, err := tc.service.Get(taskID)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	tc.printTask(task)
}

func (tc *TaskController) handleList(cmdArgs []string) {
	opts := service.TaskListOptions{}
	availableParams := []CommandParam{
		{"status", true}, {"search", true},
		{"limit", true}, {"page", true},
		{"overdue", false},
	}
	cmdParams, err := tc.parseParams(strings.Join(cmdArgs, " "), availableParams)
	if err != nil {
		fmt.Printf("Ошибка: %v", err)
		return
	}
	for key, value := range cmdParams {
		switch key {
		case "overdue":
			opts.Overdue = true
		case "status":
			status := domain.TaskStatus(value)
			opts.Status = &status
		case "limit":
			limit, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("Ошибка: invalid value of limit")
				return
			}
			opts.Limit = limit
		case "page":
			page, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("Ошибка: invalid value of page")
				return
			}
			opts.Page = page
		case "search":
			opts.Search = value
		}
	}
	listOut, err := tc.service.List(opts)
	if err != nil {
		if errors.Is(err, domain.ErrorTasksNotFound) {
			fmt.Println("Задач не найдено")
			return
		} else {
			fmt.Printf("Ошибка: %v", err)
			return
		}

	}
	tc.printList(listOut)
}

func (tc *TaskController) handleChStatus(cmdArgs []string) {
	if len(cmdArgs) < 2 {
		fmt.Println("Ошибка: недостаточно аргументов")
	}
	id, err := strconv.Atoi(cmdArgs[0])
	if err != nil {
		fmt.Println("Ошибка: некорректный ID")
		return
	}
	status := domain.TaskStatus(cmdArgs[1])
	task, err := tc.service.UpdateStatus(domain.TaskID(id), status)
	if err != nil {
		if errors.Is(err, domain.ErrorStatusNotExist) {
			fmt.Println("Ошибка: такого статуса нет")
			return
		} else {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}
	}

	tc.printTask(task)
}

func (tc *TaskController) handleChTask(cmdArgs []string) {
	if len(cmdArgs) < 2 {
		fmt.Println("Ошибка: недостаточно аргументов")
		return
	}
	availableParams := []CommandParam{
		{"title", true},
		{"dl", true},
	}
	id, err := strconv.Atoi(cmdArgs[0])
	if err != nil {
		fmt.Println("Ошибка: некорректный ID")
		return
	}
	cmdParams, err := tc.parseParams(strings.Join(cmdArgs[1:], " "), availableParams)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	updateOpts := service.TaskUpdateOptions{}
	for key, value := range cmdParams {
		switch key {
		case "title":
			updateOpts.Title = value
		case "dl":
			deadline, err := time.Parse(time.DateOnly, value)
			if err != nil {
				fmt.Printf("Ошибка: %v\n", err)
				fmt.Printf("Проверьте, что дата соответствует формату - %s\n", time.DateOnly)
				return
			}
			updateOpts.Deadline = deadline
		}
	}
	task, err := tc.service.Update(domain.TaskID(id), updateOpts)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	tc.printTask(task)
}

func (tc *TaskController) handleDelTask(cmdArgs []string) {
	if len(cmdArgs) < 1 {
		fmt.Println("Ошибка: недостаточно аргументов")
		return
	}
	id, err := strconv.Atoi(cmdArgs[0])
	if err != nil {
		fmt.Println("Ошибка: некорректный ID")
		return
	}

	err = tc.service.Delete(domain.TaskID(id))
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Println("Задача удалена.")
}

func (tc *TaskController) handleExit() {
	os.Exit(0)
}

func (tc *TaskController) printList(listOut domain.ListTaskOutput) {
	var b strings.Builder
	startNumeric := (listOut.Page-1)*listOut.Limit + 1
	fmt.Printf("%-4s %-7s %-30s %-15s %s\n", "№", "ID", "Title", "Status", "Deadline")
	for _, task := range listOut.Tasks {
		fmt.Fprintf(&b, "%-4d %-7d %-30s %-15s %s\n",
			startNumeric, task.ID,
			task.Title, task.Status,
			task.Deadline.Format(time.DateOnly),
		)
		startNumeric++
	}
	fmt.Fprintln(&b, "")
	fmt.Fprintf(&b, "Current page: %d of %d", listOut.Page, listOut.TotalPages)

	fmt.Println(b.String())
}

func (tc *TaskController) printTask(task domain.Task) {
	fmt.Println("ID:", task.ID)
	fmt.Println(task.Title)
	fmt.Println("Статус:", task.Status)
	fmt.Println("Дедлайн:", task.Deadline.Format(time.DateOnly))
}

func (tc *TaskController) parseParams(s string, availableParams []CommandParam) (map[string]string, error) {
	cmdParams := make(map[string]string)
	paramsMap := make(map[string]CommandParam)
	for _, cmd := range availableParams {
		paramsMap[cmd.Cmd] = cmd
	}
	for len(s) > 0 {
		for len(s) > 0 && s[0] == ' ' {
			s = s[1:]
		}
		if len(s) == 0 {
			break
		}

		i := 0
		for i < len(s) && s[i] != ' ' {
			i++
		}
		cmd := s[:i]
		cmdParam, exists := paramsMap[cmd]
		if !exists {
			return nil, fmt.Errorf("invalid params")
		}
		if !cmdParam.RequiredValue {
			cmdParams[cmdParam.Cmd] = ""
			continue
		}
		s = s[i:]
		if len(s) == 0 {
			break
		}
		s = s[1:]
		var value string
		i = 0
		if s[0] == '"' || s[0] == '\'' {
			quote := s[0]
			s = s[1:]
			for i < len(s) && s[i] != quote {
				i++
			}
			value = s[:i]
			if i == len(s) {
				return nil, fmt.Errorf("invalid arg for %s", cmd)
			}
			s = s[i+1:]
		} else {
			for i < len(s) && s[i] != ' ' {
				i++
			}
			value = s[:i]
			s = s[i:]
		}

		cmdParams[cmd] = value
	}

	return cmdParams, nil
}
