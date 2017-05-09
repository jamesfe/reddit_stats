package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
	"strconv"
	"strings"

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
	f, err := os.Open(inFile)
	if err != nil {
		log.Debug("open")
		log.Fatal(err)
	}
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
	var authorsPerDay, err = UniqueAuthorsPerDayAnalysis(simpleParams)
	if err == nil {
		a, b := json.Marshal(authorsPerDay)
		if b == nil {
			log.Infof("%s", a)
		} else {
			log.Errorf("Error parsing output JSON: %s", b)
		}
		log.Infof("UniqueAuthors: %#v", authorsPerDay)
	} else {
		log.Errorf("There was an error performing UniqueAuthorAnalysis: %s", err)
	}

}
