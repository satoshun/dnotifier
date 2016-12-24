package dnotifier

import "fmt"

// NewLocal create Local
func NewLocal() *Local {
	return new(Local)
}

// Local has local params
type Local struct {
}

// SendMessage send message to locals
func (s *Local) SendMessage(message Message) error {
	body := message.Diff
	fmt.Println(body)
	return nil
}
