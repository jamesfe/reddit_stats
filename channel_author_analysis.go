package main

import (
	"encoding/json"
	"strings"
	"sync"
	"time"
)

type AuthorDateTuple struct {
	AuthorDate string
	AuthorName string
}

func AuthorSingleLine(line []byte, ret chan<- AuthorDateTuple, wg *sync.WaitGroup) {
	defer wg.Done()
	if isDonaldLite(line) {
		var rawJsonMap interface{}
		jumerr := json.Unmarshal(line, &rawJsonMap)

		if jumerr == nil {
			var result AuthorDateTuple
			v := rawJsonMap.(map[string]interface{})
			subreddit := v["subreddit"].(string)
			if strings.ToLower(subreddit) == "the_donald" {
				result.AuthorName = v["author"].(string)
				realTime := getIntTimestamp(v["created_utc"])
				if realTime != 0 && result.AuthorName != "[deleted]" { // if it is junk, don't record
					result.AuthorDate = time.Unix(int64(realTime), 0).Format("02-01-2006")
					ret <- result
				}
			}
		} else {
			log.Warningf("Bad JSON Unmarshall: %s", jumerr)
		}
	}
	log.Errorf("error done")
}
