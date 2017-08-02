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
		outPerDay := CreateActiveUserMap(longevityMap, minDate, maxDate, weekLength, minSeconds, utils.GetWeekString)
		longevityOutput := AggregateByAuthorLongevity(longevityMap, minSeconds)
		utils.DumpJSONToFile("longevity", longevityOutput)
		utils.DumpJSONToFile("outPerDay", outPerDay)
	}
	if config.AnalysisConfiguration.AnalysisMap["unique_author_count"] {
		uniqueAuthorMap := AggregateUniqueAuthors(far)
		utils.DumpJSONToFile("uniqueauthors", uniqueAuthorMap)
	}
	if config.AnalysisConfiguration.AnalysisMap["deleted"] {
		outputMap := AggregateByDeletedCommentCounts(far)
		utils.DumpJSONToFile("deleted", outputMap)
	}
}

func AggregateUniqueAuthors(analysisResults map[string]map[string]int) analysis.UniqueAuthorsPerDayResult {
	resultMap := analysis.NewUniqueAuthorsPerDayResult()
	for k, v := range analysisResults {
		(*resultMap).AuthorsPerDay[k] = len(v)
	}
	return *resultMap
}

func CreateActiveUserMap(input map[string]*data_types.UserLongevityResult, start int, end int, delta int, minSecondsDiff int, dateFunc data_types.DateToString) map[string]int {
	/*
	   First we create a map with all the dates we want.  We need start, end, and delta and aggregate function.
	   Second we iterate over the list of items.  We generate a date, round it to delta and then add delta until it's larger than end.
	*/
	rv := make(map[string]int)
	// Make the original map
	for a := start; a < end; a += delta {
		rv[dateFunc(a)] = 0
	}

	// For each item, add a bunch of things.
	for _, element := range input {
		if (element.LastPost - element.FirstPost) >= minSecondsDiff {
			for b := element.FirstPost - (element.FirstPost % delta); b < element.LastPost; b += delta {
				rv[dateFunc(b)] += 1
			}
		}
	}
	return rv
}

func AggregateByAuthorLongevity(input map[string]*data_types.UserLongevityResult, minSecondsDiff int) []*data_types.TimePeriod {
	var rv []*data_types.TimePeriod
	for _, element := range input {
		if (element.LastPost - element.FirstPost) >= minSecondsDiff {
			newObject := &data_types.TimePeriod{
				StartDate: utils.GetDayString(element.FirstPost),
				EndDate:   utils.GetDayString(element.LastPost)}
			rv = append(rv, newObject)
		}
	}
	return rv
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
	}
	return outputMap
}
