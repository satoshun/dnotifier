package main

import (
	"flag"
	"fmt"
	"log"
	"path"
	"path/filepath"

	dnotifier "github.com/satoshun/dnotifier/lib"
)

var (
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
	flag.Var(&files, "f", "specified files")
	flag.Parse()

	if *slackHookURL == "" {
		log.Fatal("please specify -u option")
	}

	if len(files) == 0 {
		log.Fatal("please specify -f option")
	}

	api := dnotifier.Slack{
		HookURL:   *slackHookURL,
		Channel:   *channel,
		IconEmoji: *iconEmoji,
	}

	for i, f := range files {
		files[i], _ = filepath.Abs(f)
		fmt.Println("Watching: " + f)
	}

	w := dnotifier.Watch(files...)

	for {
		select {
		case msg := <-w.Event:
			log.Println("change:" + msg.Path)
			if msg.Diff != "" {
				name := *userName
				if name == "" {
					name = path.Base(msg.Path)
				}
				api.PostMessage(name, msg.Diff)
			}
		}
	}
}
