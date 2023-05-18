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
		if ce.Config.Directories[event.Path].LogFile != "" {
			err = ce.logCommands(errStream.String(), outStream.String(), err, ce.Config.Directories[event.Path].LogFile)
			if err != nil {
				break
			}
		} else if err != nil {
			break
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

func (ce *CommandsExecutor) logCommands(errString string, outString string, err error, logFile string) error {
	if errString != "" {
		ce.Logger.Error(logFile, errString)
	}
	if err != nil {
		ce.Logger.Error(logFile, err.Error())
		return err
	}
	if outString != "" {
		ce.Logger.Info(logFile, outString)
	}
	return nil
}
