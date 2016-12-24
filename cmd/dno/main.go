package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	dnotifier "github.com/satoshun/dnotifier"
)

var (
	command string

	// slack params
	slackHookURL = flag.String("u", "", "your slack hook url")
	channel      = flag.String("c", "#general", "your channel name. default #general")
	userName     = flag.String("n", "", "your username")
	iconEmoji    = flag.String("i", ":ghost:", "your icon emoji")
	files        arrayFlags
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
	if len(os.Args) <= 1 {
		usage()
		os.Exit(0)
	}

	command = strings.ToLower(os.Args[1])
	if command != "slack" {
		log.Fatal("only enable Slack command")
	}
	// remove subcommand
	copy(os.Args[1:], os.Args[2:])

	flag.Var(&files, "f", "specified files")
	flag.Parse()

	if len(files) == 0 {
		log.Fatal("please specify -f option")
	}

	var ms dnotifier.Messenger
	if command == "slack" {
		if *slackHookURL == "" {
			log.Fatal("necessary webhook url: -u")
		}
		if *channel == "" {
			log.Fatal("necessary channel params: -c")
		}
		ms = dnotifier.NewSlack(*slackHookURL, *channel, *iconEmoji, *userName)
	}
	if ms == nil {
		log.Fatal("illegal args")
	}

	for i, f := range files {
		files[i], _ = filepath.Abs(f)
		fmt.Printf("Watching: %s ...\n", f)
	}

	w := dnotifier.Watch(files...)
	log.Fatal(watch(ms, w))
}

func watch(m dnotifier.Messenger, w dnotifier.Watcher) error {
	for {
		select {
		case msg := <-w.Event:
			if msg.Diff != "" {
				log.Printf("changed: %s\n", msg.Path)
				e := m.SendMessage(dnotifier.Message{Diff: msg.Diff})
				if e != nil {
					return e
				}
			}
		}
	}
}

func usage() {
	flag.Usage()
}
