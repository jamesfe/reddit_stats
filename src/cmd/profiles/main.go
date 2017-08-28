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

var log = logging.MustGetLogger("profiles")
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
	var resultItem data_types.AuthorDateSubTuple // we reuse this address for results

	authorProfileMap := make(map[string]*data_types.AuthorProfile) // a map of usernames to profiles
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
			} else if analysis.AuthorSingleLineMultiNoFilter(inputBytes, &resultItem) {
				if config.AnalysisConfiguration.AnalysisMap["profile_users"] {
					/* If we analyzing users, extract some features from them. */
					if val, ok := authorProfileMap[resultItem.AuthorName]; ok {
						// Aggregate here
						if individualProfile, ok2 := val.CommentCountsBySub[resultItem.SubReddit]; ok2 {
							val.CommentCountsBySub[resultItem.SubReddit] += 1
						} else {
							val.CommentCountsBySub[resultItem.SubReddit] = 0
						}
					} else {
						authorPofileMap[resultItem.AuthorName] = data_types.AuthorProfile{CommentCountsBySub: make(map[string]int)}
						authorProfileMap[resultItem.AuthorName].CommentCountsBySub[resultItem.SubReddit] = 1
					}
				}
			}
		}
		if lines == config.MaxLines {
			log.Infof("Max lines reached")
			break
		}
	}

	if config.AnalysisConfiguration.AnalysisMap["deleted"] {
		outputMap := analysis.AggregateByDeletedCommentCounts(far)
		utils.DumpJSONToFile("deleted", outputMap)
	}
}
