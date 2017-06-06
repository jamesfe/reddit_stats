package analysis

import (
	"encoding/json"
	"github.com/jamesfe/reddit_stats/src/data_types"
	"strings"
	"time"
)

func AuthorSingleLine(line []byte, result *data_types.AuthorDateTuple) bool {
	/* Take some bytes, convert them from JSON and return the author and date as a string
	if they are valid. */

	if isDonaldLite(line) {
		var rawJsonMap interface{}
		jumerr := json.Unmarshal(line, &rawJsonMap)

		if jumerr == nil {
			v := rawJsonMap.(map[string]interface{})
			subreddit := v["subreddit"].(string)
			if strings.ToLower(subreddit) == "the_donald" {
				result.AuthorName = v["author"].(string)
				realTime := getIntTimestamp(v["created_utc"])
				if realTime != 0 && result.AuthorName != "[deleted]" { // if it is junk, don't record
					result.AuthorDate = time.Unix(int64(realTime), 0).Format("02-01-2006")
					return true
				}
			}
		} else {
			log.Warningf("Bad JSON Unmarshall: %s", jumerr)
		}
	}
	return false
}

func AggregateAuthorLine(authorTuple *data_types.AuthorDateTuple, resultMap *map[string]map[string]int) {
	/* mutate the map in place with a new result */
	if (*resultMap)[authorTuple.AuthorDate] != nil {
		(*resultMap)[authorTuple.AuthorDate][authorTuple.AuthorName] += 1
	} else {
		(*resultMap)[authorTuple.AuthorDate] = make(map[string]int)
	}
}
