package dnotifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"context"
)

const (
	hookURL = "https://api.hipchat.com/v2/room/%s/notification?auth_token=%s"
	timeOut = 5
)

func NewHipChat(roomid, token string) *HipChat {
	return &HipChat{
		RoomID: roomid,
		Token:  token,
	}
}

// HipChat has toke room id and access token.
// See: https://<your team>.hipchat.com/rooms/
type HipChat struct {
	RoomID string
	Token  string
}

// MessageData is parameter when POST to hipchat API.
type MessageData struct {
	Color         string `json:"color"`
	MessageFormat string `json:"message_format"`
	Message       string `json:"message"`
	Notification  bool   `json:"notification"`
}

// NewMessageData is the POST data which processed the data as feeling good.
func NewMessageData(msgfmt, msg string, notify bool) *MessageData {
	return &MessageData{
		Color:         "purple",
		MessageFormat: msgfmt,
		Message:       msg,
		Notification:  notify,
	}
}

// HookURL is url of API endpoint.
func (h *HipChat) HookURL() string {
	return fmt.Sprintf(hookURL, h.RoomID, h.Token)
}

// SendMessage send message to HipChat.
func (h *HipChat) SendMessage(message Message) error {
	log.Println("Send to hipchat.")

	hostName, err := os.Hostname()
	if err != nil {
		log.Println("Can't get hostname.: ", err)
		return err
	}

	msgData := NewMessageData(
		"text",
		hostName+"\n"+message.Diff,
		true,
	)

	jsonData, err := json.Marshal(msgData)
	if err != nil {
		log.Println(err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut*time.Second)
	defer cancel()
	req, err := http.NewRequest("POST", h.HookURL(), bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(resp)
		return err
	}
	log.Println("Done sending a diff log.")
	return err
}
