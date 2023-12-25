package task

import "fmt"

// Task represents a one-off task that can be executed and waited for. Tasks
// can depend on other tasks. When a task is executed it will execute all its
// dependencies first. Once a task is executed it cannot be executed again
// (calling Execute() is a no-op and Wait() will return immediately with the
// error returned by the first execution).
type Task interface {
	fmt.Stringer

	// Execute executes the task and all its dependencies. If the task has
	// already been executed, this is a no-op.
	Execute()

	// Wait waits for the task to complete and returns the error returned by
	// the task's function. If the task has already been executed, this
	// returns immediately with the error returned by the previous execution.
	Wait() error
}

// New creates a new task with the given name, function, dependencies and
// arguments. T is the type of the argument the task accepts.
func New[T any](name string, f func(T) error, deps []Task, args T) Task {
	return newTaskImpl[T](name, f, deps, args)
}
