package analysis

import (
	"testing"

	"github.com/jamesfe/reddit_stats/src/data_types"
	"github.com/jamesfe/reddit_stats/src/utils"
)

func TestAuthorSingleLine(t *testing.T) {
	data := utils.LoadRawByteArray("../../test_data/single_line_data.json", 8)

	var numDonaldRows int = 0
	var numExpectedRows int = 4
	var uniqueAuthors int = 4
	var authorList []string
	var isChecked bool = false
	var res data_types.AuthorDateTuple // we reuse this address for results
	for i := 0; i < len(data); i++ {
		isChecked = AuthorSingleLine(data[i], &res, utils.GetWeekString, false)
		if isChecked {
			numDonaldRows++
			authorList = append(authorList, res.AuthorName)
		}
	}
	if numDonaldRows != numExpectedRows {
		t.Errorf("Expected number of rows to be %d but found %d", numExpectedRows, numDonaldRows)
	}
	if len(authorList) != uniqueAuthors {
		t.Errorf("Found these authors %s but only expected %d.", authorList, uniqueAuthors)
	}
}

func TestBadAuthorSingleLine(t *testing.T) {
	var res data_types.AuthorDateTuple // we reuse this address for results
	val := AuthorSingleLine([]byte("blah t5_38unr blah"), &res, utils.GetWeekString, false)
	if val != false {
		t.Error("Bad JSON resulted in a true.")
	}
}

func TestAggregateAuthorLine(t *testing.T) {
	tmap := make(map[string]map[string]int)
	var res data_types.AuthorDateTuple = data_types.AuthorDateTuple{AuthorName: "blah", AuthorDate: "02-2017", Timestamp: 1000}
	AggregateAuthorLine(&res, &tmap)
	if tmap["02-2017"][res.AuthorName] != 1 {
		t.Errorf("Did not register first author.")
	}
	AggregateAuthorLine(&res, &tmap)
	if tmap["02-2017"][res.AuthorName] != 2 {
		t.Errorf("Did not register second author.")
	}
}
