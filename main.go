package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	dnotifier "github.com/satoshun/dnotifier/lib"
)

var (
	slackHookURL = flag.String("u", "", "your slack hook url")
	channel      = flag.String("c", "general", "your channel name. default general")
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
		HookURL: *slackHookURL,
	}

	for i, f := range files {
		files[i], _ = filepath.Abs(f)
		fmt.Println("Watching: " + f)
	}

	w := dnotifier.Watch(files...)

	for {
		select {
		case msg := <-w.Event:
			log.Println("post:" + msg)
			if msg != "" {
				api.PostMessage(msg)
			}
		}
	}
}
