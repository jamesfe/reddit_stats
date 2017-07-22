package utils

import (
	"os"
	"os/exec"
	"testing"
)

func TestGetDaystring(t *testing.T) {
	tryAndExpect := func(try int, expect string) {
		val := GetDayString(try)
		if val != expect {
			t.Errorf("Expected %s but function returned %s", expect, val)
		}
	}
	tryAndExpect(1500739788, "22-07-2017")
	tryAndExpect(1485993600, "02-02-2017")
	tryAndExpect(1483228800, "01-01-2017")
}

func TestGetWeekString(t *testing.T) {
	tryAndExpect := func(try int, expect string) {
		val := GetWeekString(try)
		if val != expect {
			t.Errorf("Expected %s but function returned %s", expect, val)
		}
	}
	tryAndExpect(1500739788, "29-2017")
	tryAndExpect(1485993600, "5-2017")
	tryAndExpect(1483228800, "52-2016")
}

func TestJsonFileReader(t *testing.T) {
	/* It should crash if the file is not found */
	if os.Getenv("BE_CRASHER") == "1" {
		var k string
		ReadJsonFile("blah", k)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestJsonFileReader")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
