package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/pprof"
	"time"

	"github.com/jamesfe/reddit_stats/data_types"
	"github.com/jamesfe/reddit_stats/src/protoanalysis"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("reddit_stats")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
)

var donaldSubreddit string = "t5_38unr"
var donaldBytes []byte = []byte(donaldSubreddit)

func main() {
	filename := flag.String("filename", "", "input filename")
	checkInterval := flag.Int("cv", 1000000, "check value")
	maxLines := flag.Int("maxlines", 0, "max lines to read")
	inputFormat := flag.String("informat", "json", "input type: json or proto")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	flag.Parse()
	var delim byte = '\n'
	if *inputFormat == "proto" {
		delim = 200
	}
	filesToCheck := getFilesToCheck(*filename)

	var lines int = 0
	var resultItem data_types.AuthorDateTuple // we reuse this address for results
	far := make(map[string]map[string]int)

	log.Infof("Entering analysis stream.")
	for _, file := range filesToCheck {
		inFileReader, f := getFileReader(file)
		defer f()
		for lines = lines; lines < *maxLines; lines++ {
			if lines%*checkInterval == 0 {
				log.Debugf("Read %d lines", lines)
			}
			var inputBytes, err = inFileReader.ReadBytes(delim)
			if err == nil { // really trying to isolate the business code right here so we can call one or two functions.
				switch *inputFormat {
				case "json":
					if AuthorSingleLine(inputBytes, &resultItem) {
						AggregateAuthorLine(&resultItem, &far)
					}
				case "proto":
					if protoanalysis.ProtoSingleLineAnalysis(inputBytes, &resultItem) {
						AggregateAuthorLine(&resultItem, &far)
					}
				}
			} else {
				log.Errorf("File Error: %s", err) // maybe we are in an IO error?
				break
			}
		}
		if lines == *maxLines {
			log.Infof("Max lines reached")
			break
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
