package dnotifier

import "testing"

func TestStoreInCache(t *testing.T) {
	cache = make(map[string]string)

	err := storeInCache("./differ.go")
	if err != nil {
		t.Errorf("%s", err)
	}
	if _, ok := cache["./differ.go"]; !ok {
		t.Error("cache is empty")
	}

	err = storeInCache("./hoge.go")
	if err == nil {
		t.Error("not throw error when file not exists")
	}
}
