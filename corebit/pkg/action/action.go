package action

import (
	"errors"
	"fmt"
	"time"
)

// ActionStatus represents the status of the action.
type ActionStatus string

const (
	// Pending status when an action is created but not executed.
	Pending ActionStatus = "PENDING"
	// InProgress status when an action is currently being executed.
	InProgress ActionStatus = "IN_PROGRESS"
	// Completed status when the action is executed successfully.
	Completed ActionStatus = "COMPLETED"
	// Failed status when the action fails during execution.
	Failed ActionStatus = "FAILED"
)

// Action represents an abstraction for any operational unit chosen by the LLM to execute.
type Action struct {
	Name        string                                                       // Name of the action
	Description string                                                       // Description of the action
	Input       map[string]interface{}                                       // Inputs required to perform the action
	Output      map[string]interface{}                                       // Output produced by the action
	Status      ActionStatus                                                 // Current status of the action
	Timestamp   time.Time                                                    // Timestamp for when the action was initiated
	Security    string                                                       // Security level for the action
	Handler     func(map[string]interface{}) (map[string]interface{}, error) // Handler function to execute the action
}

// NewAction creates a new action instance with default values.
func NewAction(name, description, security string, handler func(map[string]interface{}) (map[string]interface{}, error)) *Action {
	return &Action{
		Name:        name,
		Description: description,
		Input:       make(map[string]interface{}),
		Output:      make(map[string]interface{}),
		Status:      Pending,
		Timestamp:   time.Now(),
		Security:    security,
		Handler:     handler,
	}
}

// Execute performs the action, logs metrics, and stores results in the database.
func (a *Action) Execute() error {
	if a.Status != Pending {
		return errors.New("action already executed or in progress")
	}

	a.Status = InProgress
	fmt.Printf("Executing action: %s\n", a.Name)

	// Start tracking execution time
	start := time.Now()

	output, err := a.Handler(a.Input)
	if err != nil {
		a.Status = Failed
		fmt.Printf("Action %s failed: %v\n", a.Name, err)
		// Log failure metrics here (InfluxDB)
		return err
	}

	a.Output = output
	a.Status = Completed

	// Calculate execution duration
	duration := time.Since(start)
	fmt.Printf("Action %s completed in %v\n", a.Name, duration)

	// Log success metrics (InfluxDB)
	// Store results in the PostgreSQL database

	return nil
}

// UpdateInput allows updating the input parameters of the action.
func (a *Action) UpdateInput(input map[string]interface{}) {
	a.Input = input
}
