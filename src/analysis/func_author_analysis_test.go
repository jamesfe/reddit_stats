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

func TestAggregateAuthorLine(t *testing.T) {
	var tmap *map[string]map[string]int = make(map[string]map[string]int)
	var res data_types.AuthorDateTuple = data_types.AuthorDateTuple{AuthorName: "blah", AuthorDate: "02-2017", Timestamp: 1000}

}

//	res *data_types.AuthorDateTuple, resultMap *map[string]map[string]int) {
