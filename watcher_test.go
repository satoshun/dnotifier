package dnotifier

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestWatch(t *testing.T) {
	_, err := Watch("./hoge.go")
	if err == nil {
		t.Error("not throw error")
	}

	f, _ := ioutil.TempFile("", "./hoge.go")
	c, err := Watch(f.Name())
	if err != nil {
		t.Errorf("throw error: %s", err)
	}

	f.Write([]byte{'a', 'b', 'c', 'd'})
	item := <-c.Event
	if item.Path != f.Name() {
		t.Errorf("%s not expected path: %s", item.Path, f.Name())
	}
	if !strings.Contains(item.Diff, "abcd") {
		t.Error("weird diff")
	}
}
