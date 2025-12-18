package task

import "errors"

var ErrActivityNotFoundForTask = errors.New("activity not found for task")
var ErrTaskNotFound = errors.New("task not found")
var ErrAssignmentNotFound = errors.New("assignment not found")
