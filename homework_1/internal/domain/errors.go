package domain

import "errors"

var (
	ErrorTaskNotFound         = errors.New("task not found")
	ErrorTasksNotFound        = errors.New("tasks not found")
	ErrorTaskValidateTitle    = errors.New("invalid title")
	ErrorTaskValidateDeadline = errors.New("invalid deadline")
	ErrorTaskSameStatus       = errors.New("task already in this status")
	ErrorTaskLockedStatus     = errors.New("status cannot be changed")
	ErrorOptsValidatePage     = errors.New("invalid page")
	ErrorOptsValidateLimit    = errors.New("invalid limit")
	ErrorStatusNotExist       = errors.New("status doesn not exist")
)
