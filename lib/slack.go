package dnotifier

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	defaultChannel  = "#general"
	defaultUserName = "webhookbot"
)

// Slack has slack params
type Slack struct {
	HookURL  string
	Channel  string
	UserName string
}

// PostMessage send message
func (s *Slack) PostMessage(message string) {
	var body = []byte(fmt.Sprintf(`{"channel":"%s","username":"%s","text":"%s"}`,
		s.channel(),
		s.userName(),
		"```"+strings.Replace(message, "\"", "\\\"", -1)+"```"))
	req, _ := http.NewRequest("POST", s.HookURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(resp)
		log.Println(err)
	}
}

func (s *Slack) channel() string {
	if s.Channel == "" {
		return defaultChannel
	}
	return s.Channel
}

func (s *Slack) userName() string {
	if s.UserName == "" {
		return defaultUserName
	}
	return s.UserName
}
