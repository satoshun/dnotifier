package dnotifier

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
	"time"
)

// Watch is watch specified path
func Watch(paths ...string) Watcher {
	event := make(chan string)
	for _, p := range paths {
		register(p)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go func() {
		cmd := exec.Command("fswatch", paths...)
		var (
			out bytes.Buffer
			i   int
		)
		cmd.Stdout = &out
		cmd.Start()
		for {
			select {
			case <-time.After(time.Second * 3):
				lines := out.String()
				c := strings.Count(lines, "\n")
				if c != i {
					set := make(map[string]struct{})
					for _, line := range strings.Split(lines, "\n")[i:] {
						set[line] = struct{}{}
					}
					for p := range set {
						if p == "" {
							continue
						}

						event <- diff(p)
					}
				}
				i = c
			}
		}
	}()

	return Watcher{
		Event: event,
	}
}
