// @program:     flashbear
// @file:        worker.go
// @author:      ac
// @create:      2024-10-20 22:46
// @description:
package common

// Task represents a task to be executed
type Task struct {
	ID     int
	Action func() (interface{}, error)
}

// Worker represents a worker that executes tasks
type Worker struct {
	tasks    chan Task
	results  chan Result
	MaxIndex int
}

// Result represents the result of a task execution
type Result struct {
	ID    int
	Value interface{}
	Err   error
}

// NewWorker creates a new Worker
func NewWorker(bufferSize int) *Worker {
	return &Worker{
		tasks:   make(chan Task, bufferSize),
		results: make(chan Result, bufferSize),
	}
}

// AddTask adds a task to the worker
func (w *Worker) AddTask(f func() (interface{}, error)) {
	id := w.IncrementMaxIndex()
	w.tasks <- Task{
		ID:     id,
		Action: f,
	}
}

// Do starts the worker to execute tasks concurrently
func (w *Worker) Do() {
	close(w.tasks)
	for task := range w.tasks {
		go func(t Task) {
			value, err := t.Action()
			w.results <- Result{ID: t.ID, Value: value, Err: err}
		}(task)
	}
}

func (w *Worker) IncrementMaxIndex() int {
	w.MaxIndex = w.MaxIndex + 1
	return w.MaxIndex
}

func (w *Worker) WatchTaskError() error {
	doneNum := 0
	for result := range w.results {
		if result.Err != nil {
			return result.Err
		}
		doneNum++
		if doneNum == w.MaxIndex {
			close(w.results)
			break
		}
	}
	return nil
}
