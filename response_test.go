package scrapegoat

import (
	"testing"
)

type testItem struct {
	name  string
	price float32
}

func TestAnyItemType(t *testing.T) {
	item := &testItem{"test", 123}

	r := &Response{Item: item}
	if ti, ok := r.Item.(*testItem); ok {
		if ti.name != "test" {
			t.Errorf("Expected Item.name to be 'test' got %v", ti.name)
		}
	} else {
		t.Error("Assertion failed")
	}
}
