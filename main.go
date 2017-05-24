package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jamesfe/reddit_stats/reddit_proto"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("creepypacket")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
)

var donaldSubreddit string = "t5_38unr"
var donaldBytes []byte = []byte(donaldSubreddit)

func main() {
	filename := flag.String("filename", "", "input filename")
	checkInterval := flag.Int("cv", 1000000, "check value")
	maxLines := flag.Int("maxlines", 0, "max lines to read")
	purpose := flag.String("purpose", "simple", "purpose: simple or proto")
	outputDir := flag.String("output", "", "output directory")

	flag.Parse()
	if *purpose == "proto" && *outputDir == "" {
		log.Errorf("Must provide output directory for proto conversion.")
		os.Exit(1)
	}

	filesToCheck := getFilesToCheck(*filename)

	var lines int = 0
	var resultItem AuthorDateTuple // we reuse this address for results
	far := make(map[string]map[string]int)

	log.Infof("Entering analysis stream.")
	for _, file := range filesToCheck {
		inFileReader, f := getFileReader(file)
		defer f()
		outWriter, flushNClose := getFileWriter(file, *outputDir)
		defer flushNClose()
		for lines = lines; lines < *maxLines; lines++ {
			if lines%*checkInterval == 0 {
				log.Debugf("Read %d lines", lines)
			}
			var jsonBytes, err = inFileReader.ReadBytes('\n')
			if err == nil { // really trying to isolate the business code right here so we can call one or two functions.
				switch *purpose {
				case "simple":
					if AuthorSingleLine(jsonBytes, &resultItem) {
						AggregateAuthorLine(&resultItem, &far)
					}
				case "proto":
					data, worked := reddit_proto.ConvertLineToProto(jsonBytes)
					fmt.Printf("Read %d bytes and wrote %d bytes: saved %.2f%%\n", len(jsonBytes), len(data), float64(len(data))/float64(len(jsonBytes))*100)
					if worked {
						_, b := outWriter.Write(data)
						if b != nil {
							log.Fatal(b)
						}
					} else {
						log.Errorf("errors!")
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
