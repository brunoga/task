package task

import (
	"fmt"
	"sync"
)

type taskImpl[T any] struct {
	name string
	f    func(T) error
	deps []Task
	args T

	m        sync.Mutex
	done     chan struct{}
	err      error
	executed bool
}

var _ Task = (*taskImpl[int])(nil)

func newTaskImpl[T any](name string, f func(T) error, deps []Task,
	args T) *taskImpl[T] {
	return &taskImpl[T]{
		name: name,
		f:    f,
		deps: deps,
		args: args,
		done: make(chan struct{}),
	}
}

func (t *taskImpl[T]) Execute() {
	t.m.Lock()
	defer t.m.Unlock()

	if t.executed {
		// Execute was already called so it is either running or has finished.
		return
	}

	t.executed = true

	// Execute all dependent tasks. Trying to execute a task twice is ok.
	for _, dep := range t.deps {
		go func(d Task) {
			d.Execute()
		}(dep)
	}

	// Wait for all dependent tasks to finish.
	for _, dep := range t.deps {
		err := dep.Wait()
		if err != nil {
			t.err = err
			close(t.done)
			return
		}
	}

	// All dependant tasks have finished so we can execute this task.
	go func() {
		t.m.Lock()
		t.err = t.f(t.args)
		close(t.done)
		t.m.Unlock()
	}()
}

func (t *taskImpl[T]) Wait() error {
	<-t.done
	if t.err != nil {
		return fmt.Errorf("task %s failed: %s", t.name, t.err)
	}

	return nil
}

func (t *taskImpl[T]) String() string {
	return fmt.Sprintf("Task:%s, Executed:%t, Error:%s", t.name, t.executed,
		t.err)
}
