package task

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	counter := atomic.Int32{}

	// Create several taks with dependencies amongst themselves.
	t1 := New[time.Duration]("t1", func(d time.Duration) error {
		counter.Add(1)
		time.Sleep(d)
		return nil
	}, nil, 5*time.Millisecond)
	t2 := New[time.Duration]("t2", func(d time.Duration) error {
		counter.Add(-1)
		time.Sleep(d)
		return nil
	}, []Task{t1}, 10*time.Millisecond)
	t3 := New[time.Duration]("t3", func(d time.Duration) error {
		counter.Add(1)
		time.Sleep(d)
		return nil
	}, []Task{t1}, 150*time.Millisecond)
	t4 := New[time.Duration]("t4", func(d time.Duration) error {
		counter.Add(1)
		time.Sleep(d)
		return nil
	}, []Task{t3, t2}, 20*time.Millisecond)

	// Execute the tasks.
	t4.Execute()

	// Wait for the tasks to finish.
	err := t4.Wait()
	if err != nil {
		t.Fatalf("t4.Wait() returned error: %v", err)
	}
	if counter.Load() != 2 {
		t.Fatalf("counter = %d, want 2", counter.Load())
	}
}
