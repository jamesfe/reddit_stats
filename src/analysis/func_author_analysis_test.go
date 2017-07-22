package analysis

import (
	"testing"
)

func TestAuthorSingleLine(t *testing.T) {
	test_data := []string{"a",
		"b",
		"c"}
	if true == false {
		t.Errorf("%s", test_data)
	}
	// TODO: Test stuff
	//		t.Errorf("Expected %s but function returned %s", expect, val)
}
