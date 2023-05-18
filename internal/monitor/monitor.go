package monitor

import (
	"github.com/fsnotify/fsnotify"
	"github.com/mokan-r/golook/internal/config"
	"log"
	"strings"
)

type Monitor interface {
	matchesRegexp(eventWatcherPath string, filepath string) (match bool)
	Start()
}

type FSMonitor struct {
	Watcher   *fsnotify.Watcher
	EventChan chan EventTrigger
	Config    *config.Config
	stop      chan struct{}
}

type EventTrigger struct {
	Path    string
	Trigger string
}

func New(watcher *fsnotify.Watcher, config *config.Config, eventChan chan EventTrigger) (*FSMonitor, error) {
	return &FSMonitor{
		Watcher:   watcher,
		EventChan: eventChan,
		Config:    config,
		stop:      make(chan struct{}),
	}, nil
}

func (fm *FSMonitor) Start() {
	for path, _ := range fm.Config.Directories {
		err := fm.Watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		for {
			select {
			case event, ok := <-fm.Watcher.Events:
				if !ok {
					log.Println("Something went wrong while getting events from channel")
					continue
				}
				eventWatcherPath := getEventWatcherPath(fm.Watcher.WatchList(), event.Name)
				if event.Has(fsnotify.Write) && fm.matchesRegexp(eventWatcherPath, event.Name) {
					fm.EventChan <- EventTrigger{eventWatcherPath, event.Name}
				}
			case err, ok := <-fm.Watcher.Errors:
				if !ok {
					log.Println("Something went wrong while getting errors from channel")
				}
				log.Println("watcher error: ", err)
			case <-fm.stop:
				log.Println("monitor stopped")
				return
			}
		}
	}()
}

func getEventWatcherPath(watchList []string, eventPath string) string {
	var matchedPaths []string
	for _, wPath := range watchList {
		if strings.HasPrefix(eventPath, wPath) {
			matchedPaths = append(matchedPaths, wPath)
		}
	}

	var ret string
	for _, path := range matchedPaths {
		if len(ret) < len(path) {
			ret = path
		}
	}

	return ret
}

func (fm *FSMonitor) matchesRegexp(eventWatcherPath string, filepath string) (match bool) {
	if len(fm.Config.Directories[eventWatcherPath].IncludeRegexp) == 0 {
		match = true
	}

	for _, reg := range fm.Config.Directories[eventWatcherPath].IncludeRegexp {
		if reg.MatchString(filepath) {
			match = true
			break
		}
	}
	for _, reg := range fm.Config.Directories[eventWatcherPath].ExcludeRegexp {
		if reg.MatchString(filepath) {
			match = false
			break
		}
	}
	return match
}
