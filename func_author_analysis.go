package main

import (
	"encoding/json"
	"strings"
	"time"
)

type AuthorDateTuple struct {
	AuthorDate string
	AuthorName string
}

func AuthorSingleLine(line []byte, result *AuthorDateTuple) bool {

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

func AggregateAuthorLine(authorTuple *AuthorDateTuple, resultMap *map[string]map[string]int) {
	/* mutate the map in place with a new result */
	if (*resultMap)[authorTuple.AuthorDate] != nil {
		(*resultMap)[authorTuple.AuthorDate][authorTuple.AuthorName] += 1
	} else {
		(*resultMap)[authorTuple.AuthorDate] = make(map[string]int)
	}
}
