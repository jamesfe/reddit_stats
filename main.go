package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("creepypacket")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
)

var donaldSubreddit string = "t5_38unr"
var donaldBytes []byte = []byte(donaldSubreddit)

func main() {
	filename := flag.String("filename", "", "input filename")
	checkInterval := flag.Int("cv", 1000000, "check value")
	maxLines := flag.Int("maxlines", 0, "max lines to read")
	flag.Parse()

	filesToCheck := getFilesToCheck(*filename)

	var lines int = 0
	var resultItem AuthorDateTuple // we reuse this address for results
	far := make(map[string]map[string]int)

	log.Infof("Entering analysis stream.")
	for _, file := range filesToCheck {
		inFileReader, f := getFileReader(file)
		defer f()
		for lines = lines; lines < *maxLines; lines++ {
			if lines%*checkInterval == 0 {
				log.Debugf("Read %d lines", lines)
			}
			var jsonBytes, err = inFileReader.ReadBytes('\n')
			if err == nil {
				if AuthorSingleLine(jsonBytes, &resultItem) {
					if far[resultItem.AuthorDate] != nil {
						far[resultItem.AuthorDate][resultItem.AuthorName] += 1
					} else {
						far[resultItem.AuthorDate] = make(map[string]int)
					}
				}
			} else {
				log.Errorf("File Error: %s", err) // maybe we are in an IO error?
				break
			}
		}
	}

	apd, marshallErr := json.Marshal(far)
	if marshallErr == nil {
		outputFilename := fmt.Sprintf("./output/output_%d.json", time.Now().Unix())
		ioutil.WriteFile(outputFilename, apd, 0644)
		// log.Infof("JSON Output: %s", apd)
		log.Infof("Output written to %s", outputFilename)
	} else {
		log.Errorf("Error parsing output JSON: %s", marshallErr)
	}

}
