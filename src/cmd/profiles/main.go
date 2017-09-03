package main

/*
	Here we define a command that will create a profile of users and allow us to perform some person-by-person
	analysis of users of the various reddits.

	We require a list of users who have posted to /r/the_donald and a list of sub-reddits we are interested in
	using as features in our machine learning model.
*/

import (
	"flag"
	"github.com/jamesfe/reddit_stats/src/analysis"
	"github.com/jamesfe/reddit_stats/src/data_types"
	"github.com/jamesfe/reddit_stats/src/utils"
	"github.com/op/go-logging"
	"math/rand"
	"time"
)

var log = logging.MustGetLogger("profiles")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
)

type SubredditIntMap map[string]int

func analysisFunction(line []byte) {
	/* For each line, if they are an author we care about, add 1 an item in their map of subs */
	var resultItem data_types.AuthorDateSubTuple
	if analysis.AuthorSingleLineMultiNoFilter(line, &resultItem) {
		// We depend on golang initializing the int to it's "zero" value.
		if targetUsers[resultItem.AuthorName] {
			aggregateCounts[resultItem.AuthorName][resultItem.SubReddit] += 1
		}
	}
}

var donaldUserList map[string]bool

func addUserToMap(line []byte) {
	/* Add the users who post to the target sub to a map to be exported later. */
	var resultItem data_types.AuthorDateSubTuple
	if analysis.AuthorSingleLineMultiNoFilter(line, &resultItem) {
		donaldUserList[resultItem.AuthorName] = true
	}
}

func readFilesAndReturnAnalysis(analysisFunc func(line []byte), inDir string) {
	var delim byte = '\n'
	filesToCheck := utils.GetFilesToCheck(inDir)
	var lines int = 0

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
			} else {
				analysisFunction(inputBytes)
			}

			if lines == config.MaxLines {
				log.Infof("Max lines reached")
				break
			}
		}
	}
}

var config data_types.StatsConfiguration
var targetUsers map[string]bool
var aggregateCounts map[string]map[string]int // map of username -> list of subs -> count of comments

func init() {
	log.Info("Initializing")
	rand.Seed(time.Now().UTC().UnixNano())
	configFile := flag.String("config", "", "config file (see sample in repo)")
	flag.Parse()
	config = utils.LoadConfigurationFromFile(*configFile)
	targetUsers = make(map[string]bool)

	var userList data_types.JSONList
	utils.ReadJsonFile(config.ProfileConfiguration.UserListFile, &userList)
	/* This should never change after this. */
	targetUsers = utils.MakeExistenceMap(userList.Items)

	/* We initialize our aggregation variable to be a map of all the known users with an empty
	   subreddit->count map inside */
	aggregateCounts = make(map[string]map[string]int)
	for _, element := range userList.Items {
		aggregateCounts[element] = make(map[string]int)
	}

	// Initialize the user list so we know which users we have seen before.
	donaldUserList = make(map[string]bool)
	log.Info("Done Initializing")
}

func main() {
	if config.CpuProfile != "" {
		stopIt := utils.StartCPUProfile(config.CpuProfile)
		defer stopIt()
	}

	buildUserProfiles := true
	findTargetUser := false

	if buildUserProfiles {
		log.info("Building user profiles.")
		/* Read the list of usernames and build some user profiles. */
		readFilesAndReturnAnalysis(analysisFunction, config.DataSource)
		// This last function is going to modify the `aggregateCounts` variable
		utils.DumpJSONToFile("user_comments_by_subreddit", aggregateCounts)
	}

	if findTargetUsers {
		/* This section is for getting usernames. */
		log.info("Finding users of the target reddit.")
		readFilesAndReturnAnalysis(addUserToMap, config.ProfileConfiguration.FilteredDataSource)
		var users []string
		for key, _ := range donaldUserList {
			users = append(users, key)
		}
		var outputUsers data_types.JSONList
		outputUsers.Items = users

		utils.DumpJSONToFile("userlist", outputUsers)
		/* This ends the get-username section. */
	}

}
