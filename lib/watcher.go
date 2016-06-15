package dnotifier

// Watcher is watcher
type Watcher struct {
	Event <-chan string
}
