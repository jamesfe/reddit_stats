package utils

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/jamesfe/reddit_stats/src/data_types"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("reddit_stats_utils")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
)

func GetDayString(input int) string {
	/* Get the date in DD-MM-YYYY format */
	return time.Unix(int64(input), 0).Format("02-01-2006")
}

func GetWeekString(input int) string {
	/* Return the week of the year in format WW-YYYY. */
	year, week := time.Unix(int64(input), 0).ISOWeek()
	return fmt.Sprintf("%d-%d", week, year)
}

func ReadJsonFile(inFile string, jsontype interface{}) {
	/* Read a file into an object */
	file, e := ioutil.ReadFile(inFile)
	if e != nil {
		log.Debugf("File error: %v\n", e)
		os.Exit(1)
	}
	json.Unmarshal(file, &jsontype)
}

func GetIntTimestamp(v interface{}) int {
	/* Sometimes the timestamps we get are float64, null, or strings. Here we check. */
	var retVal int = 0
	switch v.(type) {
	case float64:
		retVal = int(v.(float64))
	case string:
		parsed, err := strconv.ParseInt(v.(string), 10, 32)
		if err != nil {
			retVal = 0
		} else {
			retVal = int(parsed)
		}
	default:
		retVal = 0
	}
	return retVal
}

var donaldSubreddit string = "t5_38unr"
var donaldBytes []byte = []byte(donaldSubreddit)

func IsDonaldLite(data []byte) bool {
	if bytes.Contains(data, donaldBytes) {
		return true
	}
	return false
}

func IsRandomSample(percentTrue float32) bool {
	/* given the percent value we are given, return true if that percent hits randomly. */
	val := rand.Float32()
	if val*100 < percentTrue {
		return true
	} else {
		return false
	}
}

func IsDonaldCertainly(comment data_types.Comment) bool {
	return strings.ToLower(comment.Subreddit) == "the_donald"
}

func GetBufioFileWriter(inspiration string, outPath string, ending string) (*bufio.Writer, func() error) {
	outName := filepath.Base(inspiration)
	var newPath string = ""
	if ending != "" {
		newPath = filepath.Join(outPath, outName+"."+ending)
	} else {
		newPath = filepath.Join(outPath, outName)
	}

	out, err := os.Create(newPath)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(out)
	return writer, func() error { writer.Flush(); out.Close(); return nil }
}

func GetFileWriter(inspiration string, outPath string, ending string) (*gzip.Writer, func() error) {
	outName := filepath.Base(inspiration)
	var newPath string = ""
	if ending != "" {
		newPath = filepath.Join(outPath, outName+"."+ending)
	} else {
		newPath = filepath.Join(outPath, outName)
	}

	out, err := os.Create(newPath)
	if err != nil {
		log.Fatal(err)
	}
	writer := gzip.NewWriter(out)
	return writer, func() error {
		log.Debugf("Closing and flushing.")
		e := writer.Flush()
		if e != nil {
			log.Errorf("Flush error: %s", e)
		}
		f := out.Close()
		if f != nil {
			log.Error("Close Error: %s", f)
		}
		return nil
	}
}

func GetFileReader(filename string) (*bufio.Reader, func() error) {
	/* Opens the file (appropriately, gz or not) and returns a reader. */
	f, err := os.Open(filename)
	if err != nil {
		log.Debug("open")
		log.Fatal(err)
	}
	retVal := bufio.NewReader(f)

	// If it ends in GZ it is a zip and we should ungzip it.
	if strings.HasSuffix(strings.ToLower(filename), "gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			log.Debug("gzip")
			log.Fatal(err)
		}
		retVal = bufio.NewReader(gr) // note this reader.
		return retVal, gr.Close
	} else {
		return retVal, f.Close
	}
}

func isEligibleFile(f string) bool {
	f = strings.ToLower(f)
	if strings.HasSuffix(f, ".gz") || strings.HasSuffix(f, ".json") {
		return true
	}
	return false
}

func GetFilesToCheck(inFile string) []string {
	var filesToCheck []string

	if len(inFile) <= 0 {
		return []string{}
	}
	if inFile[len(inFile)-1:] == "/" {
		files, dirErr := ioutil.ReadDir(inFile)
		if dirErr == nil {
			for _, file := range files {
				fName := file.Name()
				if isEligibleFile(fName) {
					filesToCheck = append(filesToCheck, path.Join(inFile, fName))
				}
			}
		} else {
			log.Errorf("Error dir: %s", inFile)
			log.Fatalf("Appeared to be a directory but had trouble: %s", dirErr)
		}
		// we need to loop over the files
	} else {
		filesToCheck = append(filesToCheck, inFile)
	}
	return filesToCheck
}

func LoadConfigurationFromFile(inFile string) data_types.StatsConfiguration {
	var returnValue data_types.StatsConfiguration
	log.Debugf("Loading config from %s", inFile)
	file, e := ioutil.ReadFile(inFile)
	if e != nil {
		log.Fatalf("Could not open config file, error: %s", e)
	}
	json.Unmarshal(file, &returnValue)

	if len(returnValue.AnalysisConfiguration.AnalysisTypes) > 0 {
		returnValue.AnalysisConfiguration.AnalysisMap = make(map[string]bool)
		for _, val := range returnValue.AnalysisConfiguration.AnalysisTypes {
			returnValue.AnalysisConfiguration.AnalysisMap[val] = true
		}
	}

	return returnValue
}

func MakeRedditMap(items []string) map[string]bool {
	retVals := make(map[string]bool)
	for _, val := range items {
		retVals[strings.ToLower(val)] = true
	}
	return retVals
}

func StartCPUProfile(outFile string) func() {
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	return pprof.StopCPUProfile
}

func DumpJSONToFile(name string, dt interface{}) {
	outData, me := json.Marshal(dt)
	if me == nil {
		outputFilename := fmt.Sprintf("./output/%s_%d.json", name, time.Now().Unix())
		ioutil.WriteFile(outputFilename, outData, 0644)
		log.Infof("%s Output written to %s", name, outputFilename)
	} else {
		log.Errorf("Error parsing output JSON: %s", me)
	}
}
