package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	hookURL = "https://api.hipchat.com/v2/room/%s/notification?auth_token=%s"
	timeOut = 5
)

func NewHipChat(roomid, token string) *HipChat {
	return &HipChat{
		RoomId: roomid,
		Token:  token,
	}
}

type HipChat struct {
	RoomId string
	Token  string
}

type MessageData struct {
	Color         string `json:"color"`
	MessageFormat string `json:"message_format"`
	Message       string `json:"message"`
	Notification  bool   `json:"notification"`
}

func NewMessageData(msgfmt, msg string, notify bool) *MessageData {
	return &MessageData{
		Color:         "purple",
		MessageFormat: msgfmt,
		Message:       msg,
		Notification:  notify,
	}
}

func (h *HipChat) HookURL() string {
	return fmt.Sprintf(hookURL, h.RoomId, h.Token)
}

func (h *HipChat) SendMessage(message Message) error {
	log.Println("Send to hipchat.")

	hostName, err := os.Hostname()
	if err != nil {
		log.Println("Can't get hostname.: ", err)
	}

	msgData := NewMessageData(
		"text",
		hostName+"\n"+message.Diff,
		true,
	)

	jsonBytes, err := json.Marshal(msgData)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest(
		"POST",
		h.HookURL(),
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * timeOut}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Println(resp)
	}
	return err
}
