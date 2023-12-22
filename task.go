package task

import "fmt"

// Task represents a one-off task that can be executed and waited for. Tasks
// can depend on other tasks. When a task is executed it will execute all its
// dependencies first. Once a task is executed it cannot be executed again
// (calling Execute() is a no-op and Wait() will return immediately with the
// error returned by the first execution).
type Task interface {
	fmt.Stringer

	Execute()
	Wait() error
}

// New creates a new task with the given name, function, dependencies and
// arguments. T is the type of the argument the task accepts.
func New[T any](name string, f func(T) error, deps []Task, args T) Task {
	return newTaskImpl[T](name, f, deps, args)
}
