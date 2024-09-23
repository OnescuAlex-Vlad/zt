package workflow

import (
	"log"
	"time"

	"github.com/robfig/cron/v3" // Use a cron package for scheduling
)

type Task interface {
	Initialize()
	Process(elements []interface{}, executor *Executor) []interface{}
	Finalize()
}

type Executor struct {
	workers int
	// Other fields to manage execution if needed
}

// Workflow structure to manage tasks
type Workflow struct {
	tasks   []Task
	batch   int
	workers int
	name    string
	stream  func([]interface{}) []interface{}
}

// NewWorkflow initializes a new Workflow
func NewWorkflow(tasks []Task, batch int, workers int, name string, stream func([]interface{}) []interface{}) *Workflow {
	if workers <= 0 {
		workers = len(tasks) // Set default workers to number of tasks
	}
	return &Workflow{
		tasks:   tasks,
		batch:   batch,
		workers: workers,
		name:    name,
		stream:  stream,
	}
}

// Call executes the workflow for input elements
func (w *Workflow) Call(elements []interface{}) <-chan interface{} {
	output := make(chan interface{})
	executor := NewExecutor(w.workers)

	go func() {
		defer close(output)

		// Run task initializers
		w.initialize()

		// Process elements with stream processor, if available
		if w.stream != nil {
			elements = w.stream(elements)
		}

		// Process elements in batches
		for _, batch := range w.chunk(elements) {
			for _, result := range w.process(batch, executor) {
				output <- result
			}
		}

		// Run task finalizers
		w.finalize()
	}()

	return output
}

// Schedule schedules a workflow using a cron expression and elements
func (w *Workflow) Schedule(cronExpr string, elements []interface{}, iterations *int) {
	c := cron.New()
	_, err := c.AddFunc(cronExpr, func() {
		for range w.Call(elements) {
			// Process the elements
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling cron job: %v", err)
	}
	c.Start()

	// Wait for the specified iterations
	if iterations != nil {
		for i := 0; i < *iterations; i++ {
			time.Sleep(time.Duration(time.Now().UnixNano()) * time.Nanosecond)
		}
		c.Stop() // Stop cron if iterations are reached
	}
}

// Initialize runs task initializer methods (if any) before processing data
func (w *Workflow) initialize() {
	for _, task := range w.tasks {
		task.Initialize()
	}
}

// Chunk splits elements into batches
func (w *Workflow) chunk(elements []interface{}) [][]interface{} {
	var batches [][]interface{}
	for i := 0; i < len(elements); i += w.batch {
		end := i + w.batch
		if end > len(elements) {
			end = len(elements)
		}
		batches = append(batches, elements[i:end])
	}
	return batches
}

// Process processes a batch of data elements
func (w *Workflow) process(elements []interface{}, executor *Executor) []interface{} {
	for _, task := range w.tasks {
		log.Printf("Running Task")
		elements = task.Process(elements, executor)
	}
	return elements
}

// Finalize runs task finalizer methods (if any) after all data processed
func (w *Workflow) finalize() {
	for _, task := range w.tasks {
		task.Finalize()
	}
}

// NewExecutor creates a new executor instance
func NewExecutor(workers int) *Executor {
	return &Executor{
		workers: workers,
	}
}