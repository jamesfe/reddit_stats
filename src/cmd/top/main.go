package main

/*
This command exists so that we can provide a list of some "top" reddits that we are interested in.

We use the `AuthorSingleLineMulti` function to return both the author and the sub which means we can then
filter things by sub.
*/

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
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
	rand.Seed(time.Now().UTC().UnixNano())

	configFile := flag.String("config", "", "config file (see sample in repo)")
	flag.Parse()
	config := utils.LoadConfigurationFromFile(*configFile)

	if config.CpuProfile != "" {
		stopIt := utils.StartCPUProfile(config.CpuProfile)
		defer stopIt()
	}

	var targetReddits data_types.JSONList
	utils.ReadJsonFile(config.FilterConfiguration.SubredditListFile, &targetReddits)
	rmap := utils.MakeRedditMap(targetReddits.Items)

	var delim byte = '\n'
	filesToCheck := utils.GetFilesToCheck(config.DataSource)

	var lines int = 0
	var resultItem data_types.AuthorDateSubTuple // we reuse this address for results

	// Represents day -> author -> posts
	far := make(map[string]map[string]map[string]int)

	log.Infof("Entering analysis loop.")
	for _, file := range filesToCheck {
		log.Debugf("Reading %s", file)
		inFileReader, f := utils.GetFileReader(file)
		defer f()
	lineloop:
		for lines = lines; lines < config.MaxLines; lines++ {
			if lines%config.CheckInterval == 0 {
				log.Debugf("Read %d lines", lines)
			}
			if inputBytes, err := inFileReader.ReadBytes(delim); err != nil {
				log.Errorf("File Error: %s", err) // maybe we are in an IO error?
				break lineloop
			} else if analysis.AuthorSingleLineMulti(inputBytes, &resultItem, utils.GetWeekString, rmap) {
				analysis.AggregateAuthorLineMulti(&resultItem, &far)
			}
		}
		if lines == config.MaxLines {
			log.Infof("Max lines reached")
			break
		}
	}

	var allResults data_types.DeletedByDateAndReddit
	allResults.Reddits = make(map[string]data_types.DeletedByDate)
	// Now we can aggregate differently.
	for key, value := range far { // key here is the reddit and the value is a map of dates
		allResults.Reddits[key] = AggregateByDeletedCommentCounts(value)
	}

	// JSON Output
	apd, marshallErr := json.Marshal(allResults)
	if marshallErr == nil {
		outputFilename := fmt.Sprintf("./output/output_%d.json", time.Now().Unix())
		ioutil.WriteFile(outputFilename, apd, 0644)
		// log.Infof("JSON Output: %s", apd)
		log.Infof("Output written to %s", outputFilename)
	} else {
		log.Errorf("Error parsing output JSON: %s", marshallErr)
	}
}

func AggregateByDeletedCommentCounts(analysisResults map[string]map[string]int) data_types.DeletedByDate {
	/* Count up the number of deleted, total, and not-deleted comments per time period and return them in a map. */
	var today_sum int
	var deleted_sum int
	var total_sum int
	outputMap := make(map[string]data_types.DeletedTuple)
	for key, element := range analysisResults {
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
	var retVals data_types.DeletedByDate
	retVals.Dates = outputMap
	return retVals
}
