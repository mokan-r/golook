package executor

import (
	"bytes"
	"github.com/mokan-r/golook/internal/config"
	"github.com/mokan-r/golook/internal/logger"
	"github.com/mokan-r/golook/internal/monitor"
	"github.com/mokan-r/golook/pkg/models"
	"os/exec"
	"strings"
	"time"
)

type Executor interface {
	Start()
	Commands() chan models.Commands
}

type CommandsExecutor struct {
	EventsChan   chan monitor.EventTrigger
	Config       *config.Config
	Logger       logger.Logger
	CommandsChan chan models.Commands
}

func New(eventsChan chan monitor.EventTrigger, cfg *config.Config) *CommandsExecutor {
	return &CommandsExecutor{
		EventsChan:   eventsChan,
		Config:       cfg,
		Logger:       &logger.StdLogger{},
		CommandsChan: make(chan models.Commands),
	}
}

func (ce *CommandsExecutor) Start() {
	go func() {
		for {
			event := <-ce.EventsChan
			go ce.runCommands(event, ce.CommandsChan)
		}
	}()
}

func (ce *CommandsExecutor) Commands() chan models.Commands {
	return ce.CommandsChan
}

func (ce *CommandsExecutor) runCommands(event monitor.EventTrigger, commandsChan chan models.Commands) {
	for _, command := range ce.Config.Directories[event.Path].Commands {
		commandSplit := strings.Split(command, " ")
		cmd := exec.Command(commandSplit[0], commandSplit[1:]...)
		cmd.Dir = event.Path
		var outStream bytes.Buffer
		var errStream bytes.Buffer
		cmd.Stdout = &outStream
		cmd.Stderr = &errStream
		startTime := time.Now()
		err := cmd.Run()
		endTime := time.Now()
		if errStream.String() != "" {
			ce.Logger.Error(ce.Config.Directories[event.Path].LogFile, errStream.String())
		}
		if err != nil {
			ce.Logger.Error(ce.Config.Directories[event.Path].LogFile, err.Error())
			break
		}
		if outStream.String() != "" {
			ce.Logger.Info(ce.Config.Directories[event.Path].LogFile, outStream.String())
		}
		commandsChan <- models.Commands{
			Path:            event.Path,
			ChangedFile:     event.Trigger,
			ExecutedCommand: command,
			StartedAt:       startTime,
			FinishedAt:      endTime,
			ExitCode:        cmd.ProcessState.ExitCode(),
		}
	}
}
