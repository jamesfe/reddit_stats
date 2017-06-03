package main

import (
	"flag"

	"github.com/jamesfe/reddit_stats/reddit_proto"
	"github.com/jamesfe/reddit_stats/src/utils"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("convert")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
	fromFile := flag.String("input", "", "input filename or directory")
	outDir := flag.String("outdir", "", "output directory")
	fromFormat := flag.String("from", "json", "origin format")
	maxLines := flag.Int("numlines", 10, "max lines to read")

	// Right now we only have two formats
	var ending string
	if *fromFormat == "json" {
		ending = "protodata"
	} else {
		ending = "json"
	}

	flag.Parse()
	var lines int = 0
	filesToCheck := utils.GetFilesToCheck(*fromFile)
	for _, file := range filesToCheck {
		inFileReader, f := utils.GetFileReader(file)

		for lines = lines; lines < *maxLines; lines++ {
			defer f()
			outWriter, flushNClose := utils.GetFileWriter(file, *outDir, ending)
			defer flushNClose()
			switch *fromFormat {
			case "json":
				var delim byte = 200
				var inputBytes, err = inFileReader.ReadBytes('\n')
				if err != nil {
					/* Here we read some lines and then convert them. */
					data, worked := reddit_proto.ConvertLineToProto(inputBytes)

					if worked {
						data = append(data, delim)
						_, b := outWriter.Write(data)
						log.Errorf("%v", data)
						if b != nil {
							log.Fatal(b)
						}
					} else {
						log.Errorf("errors writing to file!")
					}
				} else {
					log.Errorf("File Error: %s", err) // maybe we are in an IO error?
					break
				}
			case "proto":
				var inputBytes, err = inFileReader.ReadBytes(200)
				if err != nil {
					/* Here we read some lines and then convert them. */
					log.Info("Converting a line from Proto to JSON")
					log.Infof("%v", inputBytes)
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
	}
}
