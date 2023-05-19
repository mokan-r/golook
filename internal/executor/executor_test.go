package executor

import (
	"github.com/mokan-r/golook/internal/config"
	"github.com/mokan-r/golook/internal/monitor"
	"github.com/mokan-r/golook/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Mock implementation of logger.Logger interface for testing
type MockLogger struct{}

func (l *MockLogger) Info(logFile string, message string)  {}
func (l *MockLogger) Error(logFile string, message string) {}

func TestCommandsExecutor_RunCommands(t *testing.T) {
	// Mock event and configuration
	event := monitor.EventTrigger{
		Path:    "/path/to/directory",
		Trigger: "file.txt",
	}
	cfg := &config.Config{
		Directories: map[string]config.DirectoryConfig{
			"/path/to/directory": {
				Commands: []string{"echo Hello", "exit 1"},
			},
		},
	}

	// Create a CommandsExecutor instance with mock dependencies
	ce := &CommandsExecutor{
		EventsChan:   make(chan monitor.EventTrigger, 1),
		Config:       cfg,
		Logger:       &MockLogger{},
		CommandsChan: make(chan models.Commands, 1),
	}

	// Start the executor in a separate goroutine
	go ce.Start()

	// Create a goroutine to consume the CommandsChan channel
	go func() {
		for range ce.CommandsChan {
			// Consume the channel
		}
	}()

	// Send the event trigger to the executor
	ce.EventsChan <- event

	// Wait for some time to allow the executor to process the event
	time.Sleep(time.Millisecond * 100)

	// No assertion needed in this test as we are primarily checking for deadlock errors
}

func TestCommandsExecutor_LogCommands(t *testing.T) {
	// Create a CommandsExecutor instance with mock dependencies
	ce := &CommandsExecutor{
		Logger: &MockLogger{},
	}

	// Test logging with an error stream
	err := ce.logCommands("Error output", "", nil, "log.txt")
	assert.Error(t, err)

	// Test logging with an error
	err = ce.logCommands("", "", assert.AnError, "log.txt")
	assert.Error(t, err)

	// Test logging with an output stream
	err = ce.logCommands("", "Output message", nil, "log.txt")
	assert.NoError(t, err)
}

// Mock implementation of Executor interface for testing
type MockExecutor struct{}

func (e *MockExecutor) Start() {}
func (e *MockExecutor) Commands() chan models.Commands {
	return make(chan models.Commands)
}

func TestNew(t *testing.T) {
	// Mock dependencies
	eventsChan := make(chan monitor.EventTrigger)
	cfg := &config.Config{}

	// Create a new CommandsExecutor instance
	executor := New(eventsChan, cfg)

	// Assert the type of the executor
	assert.NotNil(t, executor)
}

func TestCommandsExecutor_Commands(t *testing.T) {
	// Mock dependencies
	eventsChan := make(chan monitor.EventTrigger)
	cfg := &config.Config{}

	// Create a CommandsExecutor instance
	executor := New(eventsChan, cfg)

	// Get the commands channel
	commandsChan := executor.Commands()

	// Assert the type of the commands channel
	_, ok := <-commandsChan
	assert.True(t, ok)
}
