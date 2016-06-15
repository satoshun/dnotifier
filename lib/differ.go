package dnotifier

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"syscall"
)

var (
	cache   map[string]string
	tempDir string
)

func init() {
	temp, err := ioutil.TempDir("", ".dnotifier")
	if err != nil {
		panic(err)
	}

	tempDir = temp

	cache = make(map[string]string)
}

func register(path string) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	cache[path] = string(dat)
}

func diff(path string) EventItem {
	n, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf(`diff -u <(echo -n '%s') %s`, strings.Replace(cache[path], "'", "\\'", -1), path))
	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	r := ""
	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 1 {
					r = stdout.String()
				}
			}
		}
	}

	cache[path] = string(n)
	return EventItem{Path: path, Diff: r}
}
