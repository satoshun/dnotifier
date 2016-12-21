package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	dnotifier "github.com/satoshun/dnotifier/lib"
)

var (
	command string

	// slack params
	slackHookURL = flag.String("u", "", "your slack hook url")
	channel      = flag.String("c", "#general", "your channel name. default #general")
	userName     = flag.String("n", "", "your username")
	iconEmoji    = flag.String("i", ":ghost:", "your icon emoji")

	// hipchat params
	roomid = flag.String("hc-room", "", "hipChat room id.")
	token  = flag.String("hc-token", "", "Your HipChat access token.")

	files arrayFlags
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "want to watch files"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	command = strings.ToLower(os.Args[1])

	// flag.ParseOne()
	copy(os.Args[1:], os.Args[2:])

	flag.Var(&files, "f", "specified files")
	flag.Parse()

	if len(files) == 0 {
		log.Fatal("please specify -f option")
	}

	var ms dnotifier.Messenger
	switch command {
	case "slack":
		if *slackHookURL == "" || *channel == "" {
			log.Fatal("webhook url and channel are required.: %s,%s", *slackHookURL, *channel)
		}
		ms = dnotifier.NewSlack(*slackHookURL, *channel, *iconEmoji, *userName)
	case "hipchat":
		if *roomid == "" || *token == "" {
			log.Fatal("roomid and token are required.: %s,%s", *roomid, *token)
		}
		ms = dnotifier.NewHipChat(*roomid, *token)
	}

	if ms == nil {
		log.Fatal("illegal args")
	}

	for i, f := range files {
		files[i], _ = filepath.Abs(f)
		fmt.Printf("Watching: %s ...\n", f)
	}

	w := dnotifier.Watch(files...)

	for {
		select {
		case msg := <-w.Event:
			if msg.Diff != "" {
				log.Printf("changed: %s\n", msg.Path)
				e := ms.SendMessage(dnotifier.Message{Diff: msg.Diff})
				if e != nil {
					log.Println(e)
				}
			}
		}
	}
}
