package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(filename string, msg string)
	Error(filename string, err string)
}

type StdLogger struct{}

func (l *StdLogger) Info(filename string, msg string) {
	l.writeToFile(filename, "[INFO]", msg)
}

func (l *StdLogger) Error(filename string, err string) {
	l.writeToFile(filename, "[ERROR]", err)
}

func (l *StdLogger) writeToFile(filename string, prefix string, format string) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file:", err)
		return
	}
	defer file.Close()

	logger := log.New(file, prefix, log.LstdFlags)
	logger.Printf(format)
}
