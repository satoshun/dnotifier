package lib

// Watcher is watcher
type Watcher struct {
	Event <-chan EventItem
}

// EventItem has information of Event
type EventItem struct {
	Path string
	Diff string
}
