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

func TestDiff(t *testing.T) {
	cache = make(map[string]string)
	cache["./differ_test.go"] = "hoge"

	item, _ := diff("./differ_test.go")
	if item.Diff == "" {
		t.Error("no diff")
	}
	if item.Diff[0:5] != "--- /" {
		t.Error("wrong prefix")
	}

	item, _ = diff("./differ_test.go")
	if item.Diff != "" {
		t.Error("diff")
	}
}
