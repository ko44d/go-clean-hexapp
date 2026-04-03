package task

import "errors"

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidTitle = errors.New("title must not be empty")
	ErrTitleBlank   = errors.New("title must not be blank")
	ErrTitleTooLong = errors.New("title must not exceed 200 characters")
)
