package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime/pprof"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jamesfe/reddit_stats/reddit_proto"
	"github.com/jamesfe/reddit_stats/src/analysis"
	"github.com/jamesfe/reddit_stats/src/data_types"
	"github.com/jamesfe/reddit_stats/src/protoanalysis"
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
	inputFormat := flag.String("informat", "json", "input type: json or proto")
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
	if *inputFormat == "proto" {
		delim = 200
		log.Infof("Delimiter set")
	}
	filesToCheck := utils.GetFilesToCheck(*filename)

	var lines int = 0
	var resultItem data_types.AuthorDateTuple // we reuse this address for results
	far := make(map[string]map[string]int)
	size := make([]byte, 4)

	log.Infof("Entering analysis loop.")
	for _, file := range filesToCheck {
		inFileReader, f := utils.GetFileReader(file)
		defer f()
	lineloop:
		for lines = lines; lines < *maxLines; lines++ {
			if lines%*checkInterval == 0 {
				log.Debugf("Read %d lines", lines)
			}
			switch *inputFormat {
			case "json":
				if inputBytes, err := inFileReader.ReadBytes(delim); err != nil {
					log.Errorf("File Error: %s", err) // maybe we are in an IO error?
					break lineloop
				} else if analysis.AuthorSingleLine(inputBytes, &resultItem) {
					analysis.AggregateAuthorLine(&resultItem, &far)
				}
			case "proto":
				// some hacking here to find three of our delimeters
				if _, sizeErr := io.ReadFull(inFileReader, size); sizeErr != nil {
					log.Errorf("Size Request Error: %s", sizeErr)
					break
				}
				dataSize := binary.LittleEndian.Uint32(size)
				readerBuf := make([]byte, dataSize)
				if _, err := io.ReadFull(inFileReader, readerBuf); err != nil {
					log.Errorf("Message read error: %s", err)
					break
				}
				/* Here we read some lines and then convert them. */
				comment := &reddit_proto.Comment{}
				unmarshalerr := proto.Unmarshal(readerBuf, comment)
				if unmarshalerr == nil {
					if protoanalysis.ProtoSingleLineAnalysis(readerBuf, &resultItem) {
						AggregateAuthorLine(&resultItem, &far)
					}
				} else {
					log.Errorf("Could not parse: %s", unmarshalerr)
				}
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
