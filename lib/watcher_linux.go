package dnotifier

import (
	"log"
	"path"

	"golang.org/x/exp/inotify"
)

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
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	dir := path.Dir(p)
	err = watcher.Watch(dir)
	if err != nil {
		log.Fatal(err)
	}

	register(p)

	go func() {
		defer watcher.Close()
		for {
			select {
			case ev := <-watcher.Event:
				if ev.Mask&(inotify.IN_CLOSE_WRITE) > 0 &&
					p == ev.Name {

					event <- diff(p)
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()
}
