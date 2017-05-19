package main

import (
	"encoding/json"
	"io"
	"strings"
	"time"
)

/*
We need to do some things:
1. Find all the authors who are the_donald commenters
2. For each author, build a list of the subreddits and number of comments in each
3. Classify each author as either pro or anti trump
4. Output the number of authors for and authors against on each day.

Single Pass:
1. List of authors per day in the_donald, plus number of comments per day per author


2. List of authors and their subreddits, per week
*/

func AuthorDeepAnalysis(parameters SimpleAnalysisParameter) (map[string]map[string]int, error) {

	// Author count per day
	countPerDay := make(map[string]map[string]int)

	log.Infof("Reading filename: %s", parameters.Filename)
	inFileReader, f := getFileReader(parameters.Filename)
	defer f()
	simpleRes := newSimpleAnalysisResult()

	// var secondsPerDay int = 60 * 60 * 24
	var looperr error
	for {
		var stuff, err = inFileReader.ReadBytes('\n')
		if err != nil {
			log.Warningf("%d, %d (initial, final) lines matched out of %d", simpleRes.TotalFirstMatches, simpleRes.TotalMatches, simpleRes.TotalLinesChecked)
			if err == io.EOF {
				looperr = nil
			} else {
				log.Errorf("IO Error: %s", err)
				looperr = err
			}
			break
		}
		if (parameters.CheckLines) && (simpleRes.TotalLinesChecked >= parameters.LinesToCheck) {
			log.Errorf("Max lines of %d exceeded: %d", parameters.LinesToCheck, simpleRes.TotalLinesChecked)
			looperr = nil
			break
		}
		if isDonaldLite(stuff) {
			simpleRes.TotalFirstMatches += 1
			var rawJsonMap interface{}
			newerr := json.Unmarshal(stuff, &rawJsonMap)

			if newerr == nil {
				v := rawJsonMap.(map[string]interface{})
				subreddit := v["subreddit"].(string)
				if strings.ToLower(subreddit) == "the_donald" {
					author := v["author"].(string)
					realTime := getIntTimestamp(v["created_utc"])
					if realTime == 0 { // if it is junk, go on
						continue
					}
					dateIndex := time.Unix(int64(realTime), 0).Format("02-01-2006")
					//var dateInSeconds int = realTime - (realTime % secondsPerDay)
					if countPerDay[dateIndex] == nil {
						countPerDay[dateIndex] = make(map[string]int)
					}
					countPerDay[dateIndex][author] += 1
					// Count the number of comments, per author, per day
					simpleRes.TotalMatches += 1
				}
			} else {
				log.Errorf("JSON Parsing Error: %s", newerr)
				looperr = newerr
				break
			}
		}
		if parameters.LogLineNotification && simpleRes.TotalLinesChecked%parameters.LineIntervalNotification == 0 {
			log.Debugf("Read %d lines", simpleRes.TotalLinesChecked)
		}
		simpleRes.TotalLinesChecked++
	}
	// This is where we build the list of counts
	return countPerDay, looperr
}
