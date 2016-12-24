package lib

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

// Watch is watch specified path
func Watch(path ...string) Watcher {
	event := make(chan EventItem)
	for _, p := range path {
		watch(p, event)
	}

	return Watcher{
		Event: event,
	}
}

func watch(p string, event chan<- EventItem) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.Add(p)
	if err != nil {
		log.Fatal(err)
	}

	register(p)

	go func() {
		defer watcher.Close()
		for {
			select {
			case ev := <-watcher.Events:
				// write event
				if ev.Op|fsnotify.Write > 0 {
					event <- diff(p)
				}
			case err := <-watcher.Errors:
				log.Println("error: ", err)
			}
		}
	}()
}
