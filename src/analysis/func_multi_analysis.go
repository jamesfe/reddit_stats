package analysis

import (
	"encoding/json"
	"strings"

	"github.com/jamesfe/reddit_stats/src/data_types"
	"github.com/jamesfe/reddit_stats/src/utils"
)

func AuthorSingleLineMulti(line []byte, result *data_types.AuthorDateSubTuple, dateAgg data_types.DateToString, targetReddits map[string]bool) bool {
	/* Take some bytes, convert them from JSON and return the author, subreddit, and date as a string
	if they are valid. */

	var rawJsonMap interface{}
	jumerr := json.Unmarshal(line, &rawJsonMap)

	if jumerr == nil {
		v := rawJsonMap.(map[string]interface{})
		subreddit := v["subreddit"].(string)
		subLower := strings.ToLower(subreddit)
		if targetReddits[subLower] { // By default if the reddit is not in the dict, we get false here.
			result.AuthorName = v["author"].(string)
			realTime := utils.GetIntTimestamp(v["created_utc"])
			result.SubReddit = subreddit
			if realTime != 0 { // if it is junk, don't record
				result.AuthorDate = dateAgg(realTime)
				return true
			}
		}
	} else {
		log.Warningf("Bad JSON Unmarshall: %s", jumerr)
	}
	return false
}

func AggregateAuthorLineMulti(res *data_types.AuthorDateSubTuple, resultMap *map[string]map[string]map[string]int) {
	/* We are probably keeping track of too much here, but let's try:
	- We have a map of subreddits, each of with has a map of dates, each of which has a map of users & the # of posts they have made
	*/
	if (*resultMap)[res.SubReddit] != nil {
		if (*resultMap)[res.SubReddit][res.AuthorDate] != nil {
			(*resultMap)[res.SubReddit][res.AuthorDate][res.AuthorName] += 1
		} else {
			(*resultMap)[res.SubReddit][res.AuthorDate] = make(map[string]int)
			(*resultMap)[res.SubReddit][res.AuthorDate][res.AuthorName] = 1
		}
	} else {
		// If the subreddit item does not exist, create a date and author entry
		(*resultMap)[res.SubReddit] = make(map[string]map[string]int) // make a map of dates -> users -> ints
		(*resultMap)[res.SubReddit][res.AuthorDate] = make(map[string]int)
		(*resultMap)[res.SubReddit][res.AuthorDate][res.AuthorName] = 1

	}
}
