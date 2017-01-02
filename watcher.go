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
func Watch(paths ...string) Watcher {
	event := make(chan EventItem)
	watch(paths, event)

	return Watcher{
		Event: event,
	}
}

func watch(paths []string, event chan<- EventItem) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range paths {
		err = watcher.Add(p)
		if err != nil {
			log.Fatal(err)
		}
		err = storeInCache(p)
		if err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case ev := <-watcher.Events:
				// write or rename event
				if ev.Op&fsnotify.Write > 0 || ev.Op&fsnotify.Rename > 0 {
					event <- diff(ev.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error: ", err)
			}
		}
	}()
}
