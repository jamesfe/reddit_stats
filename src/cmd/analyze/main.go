package main

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

	flag.Parse()

	if config.CpuProfile != "" {
		stopIt := utils.StartCPUProfile(config.CpuProfile)
		defer stopIt()
	}

	var delim byte = '\n'
	filesToCheck := utils.GetFilesToCheck(config.DataSource)

	var lines int = 0
	var resultItem data_types.AuthorDateTuple // we reuse this address for results

	// Represents day -> author -> posts
	far := make(map[string]map[string]int)
	longevityMap := make(map[string]data_types.UserLongevityResult)

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
			} else if analysis.AuthorSingleLine(inputBytes, &resultItem, utils.GetWeekString, false) {
				if config.AnalysisConfiguration.AnalysisMap["unique_author_count"] == true {
					analysis.AggregateAuthorLine(&resultItem, &far)
				}
				if config.AnalysisConfiguration.AnalysisMap["author_longevity"] == true {
					analysis.AggregateLongevityLine(&resultItem, &longevityMap)
				}
			}
		}
		if lines == config.MaxLines {
			log.Infof("Max lines reached")
			break
		}
	}

	// Now we can aggregate differently.
	outputMap := AggregateByDeletedCommentCounts(far)

	// JSON Output
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

func AggregateByDeletedCommentCounts(analysisResults map[string]map[string]int) map[string]data_types.DeletedTuple {
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
	return outputMap
}
