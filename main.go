package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
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
	Body             string `json:"body"`
	Controversiality int    `json:"controversiality"`
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

func isDonaldLite(data []byte) bool {
	if strings.Contains(strings.ToLower(string(data)), "t5_38unr") { //t5_2qh3l  // t5_38unr
		return true
	}
	return false
}

func isDonaldCertainly(comment Comment) bool {
	return strings.ToLower(comment.Subreddit) == "the_donald"
}

func getFileReader(filename string) (*bufio.Reader, func() error) {
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

func main() {
	filename := flag.String("filename", "", "input filename")
	sizeCheck := flag.Int("cv", 1000000, "check value")
	maxLines := flag.Int("maxlines", 0, "max lines to read")
	flag.Parse()
	inFile = *filename

	maxLinesVal := uint32(*maxLines)
	sizeCheckVal := uint32(*sizeCheck)
	log.Debug("reading " + inFile)

	inFileReader, f := getFileReader(inFile)
	defer f()

	// now we read it, line by line in json
	var totalLines uint32 = 0
	var totalMatches int = 0
	var totalFirstMatches int = 0

	for {
		var v Comment
		var stuff, err = inFileReader.ReadBytes('\n')
		// fmt.Printf("%v", stuff)
		if err != nil {
			fmt.Printf("%d, %d (initial, final) lines matched out of %d", totalFirstMatches, totalMatches, totalLines)
			log.Fatal("READLINE: ", err)
		}
		totalLines++
		if (maxLinesVal != 0) && (totalLines >= maxLinesVal) {
			log.Errorf("Max lines of %d exceeded: %d", maxLinesVal, totalLines)
			break
		}
		if totalLines%sizeCheckVal == 0 {
			log.Debugf("Read %d lines", totalLines)
		}
		// if strings.Contains(strings.ToLower(string(stuff)), "the_donald") {
		if isDonaldLite(stuff) {
			newerr := json.Unmarshal(stuff, &v)

			if newerr == nil {
				if isDonaldCertainly(v) {
					totalMatches += 1
				}
			} else {
				log.Fatal(newerr)
			}

		}
	}
	fmt.Printf("%d, %d (initial, final) lines matched out of %d", totalFirstMatches, totalMatches, totalLines)
}
