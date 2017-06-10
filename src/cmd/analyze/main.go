package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/pprof"
	"time"

	"github.com/jamesfe/reddit_stats/src/analysis"
	"github.com/jamesfe/reddit_stats/src/data_types"
	"github.com/jamesfe/reddit_stats/src/utils"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("convert")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
	filename := flag.String("filename", "", "input filename")
	checkInterval := flag.Int("cv", 1000000, "check value")
	maxLines := flag.Int("maxlines", 0, "max lines to read")
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
	filesToCheck := utils.GetFilesToCheck(*filename)

	var lines int = 0
	var resultItem data_types.AuthorDateTuple // we reuse this address for results

	// Represents day -> author -> posts
	far := make(map[string]map[string]int)

	log.Infof("Entering analysis loop.")
	for _, file := range filesToCheck {
		inFileReader, f := utils.GetFileReader(file)
		defer f()
	lineloop:
		for lines = lines; lines < *maxLines; lines++ {
			if lines%*checkInterval == 0 {
				log.Debugf("Read %d lines", lines)
			}
			if inputBytes, err := inFileReader.ReadBytes(delim); err != nil {
				log.Errorf("File Error: %s", err) // maybe we are in an IO error?
				break lineloop
			} else if analysis.AuthorSingleLine(inputBytes, &resultItem) {
				analysis.AggregateAuthorLine(&resultItem, &far)
			}
		}
		if lines == *maxLines {
			log.Infof("Max lines reached")
			break
		}
	}

	// aggregate things
	/*
		1. For each day, total the number of deleted and non-deleted
		2. Output this map to JSON
		3. Probably do it in a function, pass by reference
	*/

	// foreach is index, data

	var today_sum int
	var deleted_sum int
	var total_sum int
	outputMap := make(map[string]data_types.DeletedTuple)
	for key, element := range far {
		today_sum = 0
		deleted_sum = 0
		for author, count := range element {
			if author != "[deleted]" {
				today_sum += count
			} else {
				deleted_sum = count
			}
		}
		total_sum = today_sum + deleted_sum
		d := &data_types.DeletedTuple{TodayTotal: today_sum, Deleted: deleted_sum, Total: total_sum}
		outputMap[key] = *d
		// probably unnecessary but:
		log.Infof("Total: %d Deleted %d", today_sum, deleted_sum)
	}

	apd, marshallErr := json.Marshal(outputMap)
	if marshallErr == nil {
		outputFilename := fmt.Sprintf("./output/output_%d.json", time.Now().Unix())
		ioutil.WriteFile(outputFilename, apd, 0644)
		// log.Infof("JSON Output: %s", apd)
		log.Infof("Output written to %s", outputFilename)
	} else {
		log.Errorf("Error parsing output JSON: %s", marshallErr)
	}
}
