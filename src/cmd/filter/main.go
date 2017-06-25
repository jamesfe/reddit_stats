package main

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
	filename := flag.String("filename", "", "input filename")
	checkInterval := flag.Int("cv", 1000000, "check value")
	maxLines := flag.Int("maxlines", 0, "max lines to read")
	outputDir := flag.String("output", "./outfilter/", "output directory")

	flag.Parse()

	var delim byte = '\n'
	filesToCheck := utils.GetFilesToCheck(*filename)

	var lines int = 0
	log.Debug("%v", filesToCheck)
	log.Infof("Entering read/write loop.")
	for _, file := range filesToCheck {
		log.Debugf("Reading %s", file)
		inFileReader, f := utils.GetFileReader(file)
		outFileWriter, g := utils.GetBufioFileWriter(file, *outputDir, "")
		defer g()
		defer f()
	lineloop:
		for lines = lines; lines < *maxLines; lines++ {
			if lines%*checkInterval == 0 {
				log.Debugf("Read %d lines", lines)
			}
			if inputBytes, err := inFileReader.ReadBytes(delim); err != nil {
				log.Errorf("File Error: %s", err) // maybe we are in an IO error?
				break lineloop
			} else {
				if utils.IsDonaldLite(inputBytes) {
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
		if lines == *maxLines {
			log.Infof("Max lines reached")
			break
		}
	}
}
