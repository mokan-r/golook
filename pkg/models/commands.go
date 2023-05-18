package models

import "time"

type Commands struct {
	Path            string
	ChangedFile     string
	ExecutedCommand string
	ExitCode        int
	StartedAt       time.Time
	FinishedAt      time.Time
}
