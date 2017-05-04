package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/op/go-logging"
)
import "flag"

var log = logging.MustGetLogger("creepypacket")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
)

var (
	inFile string
)

func newSimpleAnalysisResult() SimpleAnalysisResult {
	return (SimpleAnalysisResult{TotalMatches: 0, TotalFirstMatches: 0, TotalLinesChecked: 0})
}

type SimpleAnalysisResult struct {
	TotalMatches      int
	TotalFirstMatches int
	TotalLinesChecked int
}

type SimpleAnalysisParameter struct {
	LinesToCheck             int  // Max number of lines per file to check
	CheckLines               bool // even check them?
	LineIntervalNotification int  // How many lines between print statements?
	LogLineNotification      bool // whether or not to print notifications at line vals
	Filename                 string
}

type Comment struct {
	Author string `json:"author"`
	// Author_flair_css_class string `json:"author_flair_css_class "`
	// Author_flair_text      string `json:"author_flair_text"`
	// Body             string `json:"body"`
	Controversiality int `json:"controversiality"`
	// Created_utc      string `json:"created_utc"`
	// distinguished - mosty null
	// Edited       string `json:"edited"`
	Gilded int    `json:"gilded"`
	Id     string `json:"id"`
	// Link_id      string `json:"link_id"`
	// Parent_id    string `json:"parent_id"`
	// Retrieved_on int    `json:"retrieved_on"`
	Score     int    `json:"score"`
	Subreddit string `json:"subreddit"`
	// Subreddit_id string `json:"subreddit_id"`
	// Ups          int    `json:"ups"`
}

var donaldSubreddit string = "t5_38unr"
var donaldBytes []byte = []byte(donaldSubreddit)

func isDonaldLite(data []byte) bool {
	if bytes.Contains(data, donaldBytes) {
		return true
	}
	return false
}

func isDonaldCertainly(comment Comment) bool {
	return strings.ToLower(comment.Subreddit) == "the_donald"
}

func getFileReader(filename string) (*bufio.Reader, func() error) {
	/* Opens the file (appropriately, gz or not) and returns a reader. */
	f, err := os.Open(inFile)
	if err != nil {
		log.Debug("open")
		log.Fatal(err)
	}
	// defer f.Close()
	retVal := bufio.NewReader(f)

	// If it ends in GZ it is a zip and we should ungzip it.
	if strings.HasSuffix(strings.ToLower(inFile), "gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			log.Debug("gzip")
			log.Fatal(err)
		}
		// defer gr.Close()
		retVal = bufio.NewReader(gr) // note this reader.
		return retVal, gr.Close
	} else {
		return retVal, f.Close
	}
}

func SimpleFileAnalysis(parameters SimpleAnalysisParameter) (SimpleAnalysisResult, error) {
	/* Opens the file, reads it, counts some things up and returns a set of results.  */
	inFileReader, f := getFileReader(inFile)
	defer f()
	results := newSimpleAnalysisResult()

	for {
		var v Comment
		var stuff, err = inFileReader.ReadBytes('\n')
		if err != nil {
			log.Warningf("%d, %d (initial, final) lines matched out of %d", results.TotalFirstMatches, results.TotalMatches, results.TotalLinesChecked)
			return results, err
		}
		if (parameters.CheckLines) && (results.TotalLinesChecked >= parameters.LinesToCheck) {
			log.Errorf("Max lines of %d exceeded: %d", parameters.LinesToCheck, results.TotalLinesChecked)
			return results, nil
		}
		if isDonaldLite(stuff) {
			results.TotalFirstMatches += 1
			newerr := json.Unmarshal(stuff, &v)

			if newerr == nil && isDonaldCertainly(v) {
				results.TotalMatches += 1
			} else {
				return results, newerr
			}
		}
		if parameters.LogLineNotification && results.TotalLinesChecked%parameters.LineIntervalNotification == 0 {
			log.Debugf("Read %d lines", results.TotalLinesChecked)
		}
		results.TotalLinesChecked++
	}
	return results, nil
}

type UniqueAuthorsPerDayResult struct {
	AuthorsPerDay map[time.Time]int
	StartDate     time.Time
	EndDate       time.Time
}

func New() *UniqueAuthorsPerDayResult {
	retVal := new(UniqueAuthorsPerDayResult)
	retVal.AuthorsPerDay = map[time.Time]int{}
	return retVal
}

func UniqueAuthorsPerDayAnalysis(parameters SimpleAnalysisParameter) (UniqueAuthorsPerDayResult, error) {
	results := New()
	log.Warningf("%#v", results)

	return *results, nil
}

func main() {
	filename := flag.String("filename", "", "input filename")
	sizeCheck := flag.Int("cv", 1000000, "check value")
	maxLines := flag.Int("maxlines", 0, "max lines to read")
	flag.Parse()
	inFile = *filename

	log.Debug("reading " + inFile)

	var simpleParams SimpleAnalysisParameter
	simpleParams.Filename = inFile
	simpleParams.LinesToCheck = *maxLines
	if simpleParams.LinesToCheck == 0 {
		simpleParams.CheckLines = false
	} else {
		simpleParams.CheckLines = true
	}
	simpleParams.LineIntervalNotification = *sizeCheck
	if simpleParams.LineIntervalNotification == -1 {
		simpleParams.LogLineNotification = false
	} else {
		simpleParams.LogLineNotification = true
	}
	/*
		var simpleResults, err = SimpleFileAnalysis(simpleParams)
		if err == nil {
			log.Infof("%#v", simpleResults)
		}
	*/

	var authorsPerDay, err = UniqueAuthorsPerDayAnalysis(simpleParams)
	if err == nil {
		log.Infof("UniqueAuthors: %#v", authorsPerDay)
	}

}
