package main

import (
	"flag"
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
	log.Debugf("Config: %$v", config)
	var lines int = 0
	var resultItem data_types.AuthorDateTuple // we reuse this address for results

	// Represents day -> author -> posts
	far := make(map[string]map[string]int)
	longevityMap := make(map[string]*data_types.UserLongevityResult)
	var minDate int = 1501100780 // a long time from now
	var maxDate int = 0
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
			} else if analysis.AuthorSingleLine(inputBytes, &resultItem, utils.GetDayString, config.AnalysisConfiguration.TargetConfig.RandomSample) {
				if config.AnalysisConfiguration.AnalysisMap["unique_author_count"] || config.AnalysisConfiguration.AnalysisMap["deleted"] {
					analysis.AggregateAuthorLine(&resultItem, &far)
				}
				if config.AnalysisConfiguration.AnalysisMap["author_longevity"] {
					if resultItem.Timestamp < minDate {
						minDate = resultItem.Timestamp
					} else if resultItem.Timestamp > maxDate {
						maxDate = resultItem.Timestamp
					}
					analysis.AggregateLongevityLine(&resultItem, &longevityMap)
				}
			}
		}
		if lines == config.MaxLines {
			log.Infof("Max lines reached")
			break
		}
	}

	if config.AnalysisConfiguration.AnalysisMap["author_longevity"] {
		weekLength := 604800
		minSeconds := config.AnalysisConfiguration.LongevityConfiguration.MinDays * 24 * 3600
		outPerDay := analysis.CreateActiveUserMap(longevityMap, minDate, maxDate, weekLength, minSeconds, utils.GetWeekString)
		longevityOutput := analysis.AggregateByAuthorLongevity(longevityMap, minSeconds)
		utils.DumpJSONToFile("longevity", longevityOutput)
		utils.DumpJSONToFile("outPerDay", outPerDay)
	}
	if config.AnalysisConfiguration.AnalysisMap["unique_author_count"] {
		uniqueAuthorMap := analysis.AggregateUniqueAuthors(far)
		utils.DumpJSONToFile("uniqueauthors", uniqueAuthorMap)
	}
	if config.AnalysisConfiguration.AnalysisMap["deleted"] {
		outputMap := analysis.AggregateByDeletedCommentCounts(far)
		utils.DumpJSONToFile("deleted", outputMap)
	}
}
