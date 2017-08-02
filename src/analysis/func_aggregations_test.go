package analysis

import (
	"strconv"
	"testing"
)

func TestAggregateUniqueAuthors(t *testing.T) {
	/* Test that the number of unique authors per day is X */
	var sampleFar map[string]map[string]int = make(map[string]map[string]int)
	t1 := "01-02-2017"
	sampleFar[t1] = make(map[string]int)
	for i := 0; i < 10; i++ {
		sampleFar[t1][strconv.Itoa(i)] = i // create an author with a number of posts
	}
	results := AggregateUniqueAuthors(sampleFar)
	if results.AuthorsPerDay[t1] != 10 {
		t.Errorf("Did not get 10 authors in AggregateUniqueAuthors")
	}
}

func TestAggregateUniqueAuthorsInTwoDays(t *testing.T) {
	/* Test that the number of unique authors per day is X */
	var sampleFar map[string]map[string]int = make(map[string]map[string]int)
	t1 := "01-02-2017"
	sampleFar[t1] = make(map[string]int)
	for i := 0; i < 10; i++ {
		sampleFar[t1][strconv.Itoa(i)] = i // create an author with a number of posts
	}
	t2 := "11-11-2015"
	sampleFar[t2] = make(map[string]int)
	for i := 0; i < 50; i++ {
		sampleFar[t2][strconv.Itoa(i)] = i // create an author with a number of posts
	}

	results := AggregateUniqueAuthors(sampleFar)
	if results.AuthorsPerDay[t1] != 10 {
		t.Errorf("Did not get 10 authors in AggregateUniqueAuthors")
	}
	if results.AuthorsPerDay[t2] != 50 {
		t.Errorf("Did not get 10 authors in two days in AggregateUniqueAuthors")
	}
}
