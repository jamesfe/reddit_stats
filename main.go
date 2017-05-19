package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/op/go-logging"
)
import "flag"

var log = logging.MustGetLogger("creepypacket")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
)

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

func readResultChan(input chan AuthorDateTuple) {
	item := <-input
	log.Infof("Got one! %v", item)
}

func isDonaldCertainly(comment Comment) bool {
	return strings.ToLower(comment.Subreddit) == "the_donald"
}

func getIntTimestamp(v interface{}) int {
	/* Sometimes the timestamps we get are float64, null, or strings. Here we check. */
	var retVal int = 0
	switch v.(type) {
	case float64:
		retVal = int(v.(float64))
	case string:
		parsed, err := strconv.ParseInt(v.(string), 10, 32)
		if err != nil {
			retVal = 0
		} else {
			retVal = int(parsed)
		}
	default:
		retVal = 0
	}
	return retVal
}

func getFileReader(filename string) (*bufio.Reader, func() error) {
	/* Opens the file (appropriately, gz or not) and returns a reader. */
	f, err := os.Open(filename)
	if err != nil {
		log.Debug("open")
		log.Fatal(err)
	}
	retVal := bufio.NewReader(f)

	// If it ends in GZ it is a zip and we should ungzip it.
	if strings.HasSuffix(strings.ToLower(filename), "gz") {
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

func NewSimpleAnalysisParameter(filename string, maxLines int, checkInterval int) SimpleAnalysisParameter {
	var simpleParams SimpleAnalysisParameter
	simpleParams.Filename = filename
	simpleParams.LinesToCheck = maxLines
	if simpleParams.LinesToCheck == 0 {
		simpleParams.CheckLines = false
	} else {
		simpleParams.CheckLines = true
	}
	simpleParams.LineIntervalNotification = checkInterval
	if simpleParams.LineIntervalNotification == -1 {
		simpleParams.LogLineNotification = false
	} else {
		simpleParams.LogLineNotification = true
	}
	return simpleParams
}

func isEligibleFile(f string) bool {
	f = strings.ToLower(f)
	if strings.HasSuffix(f, ".gz") || strings.HasSuffix(f, ".json") {
		return true
	}
	return false
}

func main() {
	filename := flag.String("filename", "", "input filename")
	// checkInterval := flag.Int("cv", 1000000, "check value")
	maxLines := flag.Int("maxlines", 0, "max lines to read")
	flag.Parse()
	inFile := *filename

	log.Debug("reading " + inFile)

	var filesToCheck []string

	if inFile[len(inFile)-1:] == "/" {
		files, dirErr := ioutil.ReadDir(inFile)
		if dirErr == nil {
			for _, file := range files {
				fName := file.Name()
				if isEligibleFile(fName) {
					filesToCheck = append(filesToCheck, path.Join(inFile, fName))
				}
			}
		} else {
			log.Fatalf("Appeared to be a directory but had trouble: %s", dirErr)
		}
		// we need to loop over the files
	} else {
		filesToCheck = append(filesToCheck, *filename)
	}
	fullAuthorResult := make(map[string]map[string]int)
	simpleResultChan := make(chan AuthorDateTuple)
	log.Infof("Entering analysis stream.")
	for _, file := range filesToCheck {

		inFileReader, f := getFileReader(file)
		defer f()
		for {
			var stuff, err = inFileReader.ReadBytes('\n')
			if err == nil {
				go AuthorSingleLine(stuff, simpleResultChan)
			} else {
				log.Errorf("File Error: %s", err)
				break
			}
		}
		go readResultChan(simpleResultChan)
		/*
			// we need to abstract this and turn it into a stream of data
			log.Infof("Changing %s", file)
			var simpleParams SimpleAnalysisParameter = NewSimpleAnalysisParameter(file, *maxLines, *checkInterval)
			authorAnalysis, deeperr := AuthorDeepAnalysis(simpleParams)
			if deeperr == nil {
				for key, value := range authorAnalysis {
					fullAuthorResult[key] = value
				}
				log.Infof("Deep analysis first step complete")
			} */
		/*
			var authorsPerDay, err = UniqueAuthorsPerDayAnalysis(simpleParams)
			if err == nil {
				for key, newValue := range authorsPerDay.AuthorsPerDay {
					if value, ok := mainResult.AuthorsPerDay[key]; ok {
						mainResult.AuthorsPerDay[key] = value + newValue
					} else {
						mainResult.AuthorsPerDay[key] = newValue
					}
				}
				log.Infof("UniqueAuthors: %#v", authorsPerDay)
			} else {
				log.Errorf("There was an error performing UniqueAuthorAnalysis: %s", err)
			}
		*/

	}
	apd, marshallErr := json.Marshal(fullAuthorResult)
	if marshallErr == nil {
		outputFilename := fmt.Sprintf("./output/output_%d.json", time.Now().Unix())
		ioutil.WriteFile(outputFilename, apd, 0644)
		log.Infof("Output written to %s", outputFilename)
	} else {
		log.Errorf("Error parsing output JSON: %s", marshallErr)
	}

}
