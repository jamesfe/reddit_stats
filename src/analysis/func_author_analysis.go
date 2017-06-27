package analysis

import (
	"encoding/json"
	"time"

	"github.com/jamesfe/reddit_stats/src/data_types"
	"github.com/jamesfe/reddit_stats/src/utils"
)

func GetDayString(input int) string {
	return time.Unix(int64(input), 0).Format("02-01-2006")
}

func AuthorSingleLine(line []byte, result *data_types.AuthorDateTuple) bool {
	/* Take some bytes, convert them from JSON and return the author and date as a string
	if they are valid. */

	// if utils.IsDonaldLite(line) {
	if utils.IsRandomSample(0.25) {
		var rawJsonMap interface{}
		jumerr := json.Unmarshal(line, &rawJsonMap)

		if jumerr == nil {
			v := rawJsonMap.(map[string]interface{})
			// subreddit := v["subreddit"].(string)
			// if strings.ToLower(subreddit) == "the_donald" {
			result.AuthorName = v["author"].(string)
			realTime := utils.GetIntTimestamp(v["created_utc"])
			if realTime != 0 { // if it is junk, don't record
				result.AuthorDate = GetDayString(realTime)
				return true
			}
			// 	}
		} else {
			log.Warningf("Bad JSON Unmarshall: %s", jumerr)
		}
	}
	return false
}

func AggregateAuthorLine(authorTuple *data_types.AuthorDateTuple, resultMap *map[string]map[string]int) {
	/* Increment the author's name in the date map by one for each comment he/she has made. */
	if (*resultMap)[authorTuple.AuthorDate] != nil {
		(*resultMap)[authorTuple.AuthorDate][authorTuple.AuthorName] += 1
	} else {
		(*resultMap)[authorTuple.AuthorDate] = make(map[string]int)
	}
}
