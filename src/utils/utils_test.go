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

func TestJsonFileReaderCrashes(t *testing.T) {
	/* It should crash if the file is not found */
	if os.Getenv("BE_CRASHER") == "1" {
		var k string
		ReadJsonFile("blah", k)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestJsonFileReaderCrashes")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestIsDonaldLite(t *testing.T) {
	true_strings := []string{"t5_38unr", "aaat5_38unr", "t5_38unraaa", "aaat5_38unraaa"}
	for b := range true_strings {
		if !IsDonaldLite([]byte(true_strings[b])) {
			t.Errorf("Expected true for %s", true_strings[b])
		}
	}
	false_strings := []string{"aaa", "", "lksjfslkdjfslkdfj"}
	for b := range false_strings {
		if IsDonaldLite([]byte(false_strings[b])) {
			t.Errorf("Expected false for %s", false_strings[b])
		}
	}
}

func TestIsEligibleFile(t *testing.T) {
	true_strings := []string{"blah.gz", "blah.GZ", "thing.gz.json", "blah.JSON"}
	false_strings := []string{"this.zip", "this.jsogz", "thisjsonzip", "blah"}
	for b := range false_strings {
		if isEligibleFile(false_strings[b]) {
			t.Errorf("Expected false for %s", false_strings[b])
		}
	}
	for b := range true_strings {
		if !isEligibleFile(true_strings[b]) {
			t.Errorf("Expected false for %s", true_strings[b])
		}
	}
}

func TestMakeRedditMap(t *testing.T) {
	/* Test that it turns a list into a dict of trues */
	items := []string{"a", "b", "c"}
	testMap := MakeRedditMap(items)
	for a := range items {
		if !testMap[items[a]] {
			t.Errorf("Did not receive a true value for %s", items[a])
		}
	}
	invalidItems := []string{"blah", "bb", "cc", ""}
	for a := range invalidItems {
		if testMap[invalidItems[a]] {
			t.Errorf("Did not receive a false value for %s", invalidItems[a])
		}
	}
}

func TestLoadConfigurationFromFile(t *testing.T) {
	/* Test we are loading proper config */
	config := LoadConfigurationFromFile("../../configs/sample.json")
	// This is ugly but at least we will get an error if it fails
	if config.DataSource != "path/to/files/or/just/pathname" {
		t.Errorf("Bad data source")
	}
	if config.CheckInterval != 1234 {
		t.Errorf("Bad check interval")
	}
	if config.MaxLines != 12345 {
		t.Errorf("Bad max interval")
	}
	if config.FilterConfiguration.SubredditListFile != "./meta/political_subreddits_jul_2017.json" {
		t.Errorf("Bad Filter List File")
	}
	if config.InputFilterConfiguration.OutputDirectory != "./filters" {
		t.Errorf("Input Filter Config, bad output dir")
	}
}
