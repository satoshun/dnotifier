package lib

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Messenger interface {
	PostMessage(message string) error
}

func NewSlack(url, channel, icon, username string) *Slack {
	return &Slack{
		HookURL:   url,
		Channel:   channel,
		IconEmoji: icon,
		UserName:  username,
	}
}

// Slack has slack params
type Slack struct {
	HookURL   string
	Channel   string
	IconEmoji string
	UserName  string
}

// PostMessage send message to slack
func (s *Slack) PostMessage(message string) error {
	body := []byte(fmt.Sprintf(`{"channel":"%s","username":"%s","icon_emoji":"%s","text":"%s"}`,
		s.Channel,
		s.UserName,
		s.IconEmoji,
		"```"+standardMessage(message)+"```"))
	req, err := http.NewRequest("POST", s.HookURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(resp)
	}
	return err
}

func (s *Slack) username() string {
	if s.UserName == "" {
		return "unknown"
	}
	return s.UserName
}

func standardMessage(message string) string {
	s := strings.Replace(message, `\n`, `\\n`, -1)
	s = strings.Replace(s, `"`, `\"`, -1)
	fmt.Println(s)
	return s
}
