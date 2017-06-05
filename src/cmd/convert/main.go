package main

import (
	"encoding/binary"
	"flag"
	"io"

	"github.com/golang/protobuf/proto"
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
		ending = "protodata.gz"
	} else {
		ending = "json.gz"
	}

	flag.Parse()
	var lines int = 0
	size := make([]byte, 4)

	filesToCheck := utils.GetFilesToCheck(*fromFile)
	for _, file := range filesToCheck {
		inFileReader, f := utils.GetFileReader(file)
		defer f()
		outWriter, flushNClose := utils.GetFileWriter(file, *outDir, ending)
		defer flushNClose()
	lineloop:
		for lines = lines; lines < *maxLines; lines++ {
			switch *fromFormat {
			case "json":
				var inputBytes, err = inFileReader.ReadBytes('\n')
				if err == nil {
					/* Here we read some lines and then convert them. */
					data, worked := reddit_proto.ConvertLineToProto(inputBytes)
					if worked {
						// We put the size of the message here and then read it specially later.
						size := make([]byte, 4)
						binary.LittleEndian.PutUint32(size, uint32(len(data)))
						data = append(size, data...)
						_, b := outWriter.Write(data)
						if b != nil {
							log.Fatal(b)
						}
					} else {
						log.Errorf("errors writing to file!")
					}
				} else {
					if err != io.EOF {
						log.Errorf("File Error: %s", err)
					}
					break lineloop
				}
			case "proto":
				/* This is more difficult, we have to read a big proto object */
				_, sizeErr := io.ReadFull(inFileReader, size)
				if sizeErr == nil {
					dataSize := binary.LittleEndian.Uint32(size)
					readerBuf := make([]byte, dataSize)
					_, err := io.ReadFull(inFileReader, readerBuf) // put the message of size X in the buf
					if err == nil {
						/* Here we read some lines and then convert them. */
						comment := &reddit_proto.Comment{}
						unmarshalerr := proto.Unmarshal(readerBuf, comment)
						if unmarshalerr == nil {
							// Success, we did it
						} else {
							log.Errorf("Could not parse: %s", unmarshalerr)
						}
					} else {
						log.Errorf("Message read error: %s", err)
					}
				} else {
					log.Errorf("Size Request Error: %s", sizeErr)
					// there has been an error reading the size
				}
			}
			if lines == *maxLines {
				log.Infof("Max lines reached")
				break
			}

		}
	}
}
