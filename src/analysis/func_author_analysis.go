package analysis

import (
	"encoding/json"
	"strings"

	"github.com/jamesfe/reddit_stats/src/data_types"
	"github.com/jamesfe/reddit_stats/src/utils"
)

func AuthorSingleLine(line []byte, result *data_types.AuthorDateTuple, dateAgg data_types.DateToString, sampling bool) bool {
	/* Take some bytes, convert them from JSON and return the author and date as a string
	if they are valid. */

	if (!sampling && utils.IsDonaldLite(line)) || utils.IsRandomSample(0.25) {
		var rawJsonMap interface{}
		jumerr := json.Unmarshal(line, &rawJsonMap)

		if jumerr == nil {
			v := rawJsonMap.(map[string]interface{})
			subreddit := v["subreddit"].(string)
			if strings.ToLower(subreddit) == "the_donald" {
				result.AuthorName = v["author"].(string)
				realTime := utils.GetIntTimestamp(v["created_utc"])
				if realTime != 0 { // if it is junk, don't record
					result.AuthorDate = dateAgg(realTime)
					result.Timestamp = realTime
					return true
				}
			}
		} else {
			log.Warningf("Bad JSON Unmarshall: %s", jumerr)
		}
	}
	return false
}

func AggregateAuthorLine(res *data_types.AuthorDateTuple, resultMap *map[string]map[string]int) {
	/* Increment the author's name in the date map by one for each comment he/she has made. */
	if (*resultMap)[res.AuthorDate] != nil {
		(*resultMap)[res.AuthorDate][res.AuthorName] += 1
	} else {
		(*resultMap)[res.AuthorDate] = make(map[string]int)
		(*resultMap)[res.AuthorDate][res.AuthorName] = 1
	}
}
