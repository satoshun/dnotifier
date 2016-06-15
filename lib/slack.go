package dnotifier

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Slack has slack params
type Slack struct {
	HookURL   string
	Channel   string
	IconEmoji string
}

// PostMessage send message
func (s *Slack) PostMessage(username, message string) {
	var body = []byte(fmt.Sprintf(`{"channel":"%s","username":"%s","icon_emoji":"%s","text":"%s"}`,
		s.Channel,
		username,
		s.IconEmoji,
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
