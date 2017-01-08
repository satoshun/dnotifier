package dnotifier

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

// Watcher is watcher
type Watcher struct {
	Event <-chan EventItem
}

// EventItem has information of Event
type EventItem struct {
	Path string
	Diff string
}

// Watch is watch specified paths
func Watch(paths ...string) (*Watcher, error) {
	event := make(chan EventItem)
	err := watch(paths, event)
	if err != nil {
		return nil, err
	}

	return &Watcher{
		Event: event,
	}, nil
}

func watch(paths []string, event chan<- EventItem) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	for _, p := range paths {
		err = watcher.Add(p)
		if err != nil {
			return err
		}
		err = storeInCache(p)
		if err != nil {
			return err
		}
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case ev := <-watcher.Events:
				// write or rename event
				if ev.Op&fsnotify.Write > 0 || ev.Op&fsnotify.Rename > 0 {
					e, err := diff(ev.Name)
					if err != nil {
						log.Fatal(err, ev)
					}
					event <- *e
				}
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()
	return nil
}
