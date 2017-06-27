package analysis

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jamesfe/reddit_stats/src/data_types"
	"github.com/jamesfe/reddit_stats/src/utils"
)

func GetDayString(input int) string {
	/* Get the date in DD-MM-YYYY format */
	return time.Unix(int64(input), 0).Format("02-01-2006")
}

func GetWeekString(input int) string {
	/* Return the week of the year in format WW-YYYY. */
	year, week := time.Unix(int64(input), 0).ISOWeek()
	return fmt.Sprintf("%d-%d", week, year)
}

type DateToString func(int) string

func AuthorSingleLine(line []byte, result *data_types.AuthorDateTuple, dateAgg DateToString) bool {
	/* Take some bytes, convert them from JSON and return the author and date as a string
	if they are valid. */

	if utils.IsDonaldLite(line) {
		//	if utils.IsRandomSample(0.25) {
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
	/* Increment the author's name in the date map by one for each comment he/she has made. */
	if (*resultMap)[authorTuple.AuthorDate] != nil {
		(*resultMap)[authorTuple.AuthorDate][authorTuple.AuthorName] += 1
	} else {
		(*resultMap)[authorTuple.AuthorDate] = make(map[string]int)
	}
}
