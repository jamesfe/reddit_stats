package main

import (
	"encoding/json"
	"io"
	"strings"
	"time"
)

type UniqueAuthorsPerDayResult struct {
	AuthorsPerDay map[string]int `json:"authors_per_day"`
	StartDate     time.Time      `json:"-"`
	EndDate       time.Time      `json:"-"`
}

func New() *UniqueAuthorsPerDayResult {
	retVal := new(UniqueAuthorsPerDayResult)
	retVal.AuthorsPerDay = map[string]int{}
	return retVal
}

func UniqueAuthorsPerDayAnalysis(parameters SimpleAnalysisParameter) (UniqueAuthorsPerDayResult, error) {
	results := New()
	// TODO: add start and end date to this analysis
	log.Infof("Reading filename: %s", parameters.Filename)
	inFileReader, f := getFileReader(parameters.Filename)
	defer f()
	simpleRes := newSimpleAnalysisResult()

	AuthPerDay := map[int]map[string]bool{}
	var secondsPerDay int = 60 * 60 * 24
	var looperr error
	for {
		var stuff, err = inFileReader.ReadBytes('\n')
		if err != nil {
			log.Warningf("%d, %d (initial, final) lines matched out of %d", simpleRes.TotalFirstMatches, simpleRes.TotalMatches, simpleRes.TotalLinesChecked)
			if err == io.EOF {
				looperr = nil
			} else {
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
				/* we do some analysis here:
				we calculate the day on which the post was made
				we check if the author is in our list of authors for that day
				if he is, we do nothing
				otherwise we add him to the list
				*/

				v := rawJsonMap.(map[string]interface{})
				subreddit := v["subreddit"].(string)
				if strings.ToLower(subreddit) == "the_donald" {
					author := v["author"].(string)
					realTime := getIntTimestamp(v["created_utc"])
					if realTime == 0 { // if it is junk, go on
						continue
					}
					var extraSeconds int = realTime % secondsPerDay
					var dateInSeconds int = realTime - extraSeconds // should be the same number for every day

					if AuthPerDay[dateInSeconds] == nil {
						AuthPerDay[dateInSeconds] = make(map[string]bool)
					}
					if !AuthPerDay[dateInSeconds][author] {
						AuthPerDay[dateInSeconds][author] = true
					}
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
	// Go through the list of Authors per day and count the size of each of these things.
	log.Infof("%#v", simpleRes)
	if looperr == nil {
		for k, v := range AuthPerDay {
			item := time.Unix(int64(k), 0).Format("02-01-2006")
			results.AuthorsPerDay[item] = len(v)
		}
	} else {
		return *results, looperr
	}
	return *results, nil
}
