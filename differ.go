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

const bashPath = "/bin/bash"

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

func storeInCache(path string) error {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	cache[path] = string(dat)
	return nil
}

func diff(path string) EventItem {
	n, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(bashPath, "-c", fmt.Sprintf(`diff -u <(echo -n '%s') %s`, strings.Replace(cache[path], "'", "\\'", -1), path))
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
