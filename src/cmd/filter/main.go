package main

/*
	Read a set of files and if the criteria match, output the line to an output directory.
	Output folder is 'input_filter_settings.output_dir' in config.json
*/

import (
	"encoding/json"
	"flag"
	"strings"

	"github.com/jamesfe/reddit_stats/src/utils"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("filter")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
	configFile := flag.String("config", "", "config file (see sample in repo)")
	flag.Parse()
	config := utils.LoadConfigurationFromFile(*configFile)

	var delim byte = '\n'
	filesToCheck := utils.GetFilesToCheck(config.DataSource)

	var lines int = 0
	log.Debug("%v", filesToCheck)
	log.Infof("Entering read/write loop.")
	for _, file := range filesToCheck {
		log.Debugf("Reading %s", file)
		inFileReader, closeInFileReader := utils.GetFileReader(file)
		outFileWriter, closeOutFileWriter := utils.GetBufioFileWriter(file, config.InputFilterConfiguration.OutputDirectory, "json")
		defer closeOutFileWriter()
		defer closeInFileReader()
	lineloop:
		for lines = lines; lines < config.MaxLines; lines++ {
			if lines%config.CheckInterval == 0 {
				log.Debugf("Read %d lines", lines)
			}
			if inputBytes, err := inFileReader.ReadBytes(delim); err != nil {
				log.Errorf("File Error: %s", err) // maybe we are in an IO error?
				break lineloop
			} else {
				if utils.IsDonaldLite(inputBytes) { // the meat & potatoes: if it is Donald, check it and write it.
					var rawJsonMap interface{}
					jumerr := json.Unmarshal(inputBytes, &rawJsonMap)

					if jumerr == nil {
						v := rawJsonMap.(map[string]interface{})
						subreddit := v["subreddit"].(string)
						if strings.ToLower(subreddit) == "the_donald" {
							outFileWriter.Write(inputBytes)
						}
					} else {
						log.Warningf("Bad JSON Unmarshall: %s", jumerr)
					}
				}
			}
		}
		if lines == config.MaxLines {
			log.Infof("Max lines reached")
			break
		}
	}
}
