package dnotifier

// Messenger is send message
type Messenger interface {
	SendMessage(message Message) error
}

// Message represent message pojo
type Message struct {
	Diff string
}
