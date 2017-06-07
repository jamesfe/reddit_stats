package utils

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jamesfe/reddit_stats/src/data_types"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("reddit_stats_utils")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
)

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

func IsDonaldCertainly(comment data_types.Comment) bool {
	return strings.ToLower(comment.Subreddit) == "the_donald"
}

func GetBufioFileWriter(inspiration string, outPath string, ending string) (*bufio.Writer, func() error) {
	outName := filepath.Base(inspiration)
	newPath := filepath.Join(outPath, outName+"."+ending)

	out, err := os.Create(newPath)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(out)
	return writer, func() error { writer.Flush(); out.Close(); return nil }

}

func GetFileWriter(inspiration string, outPath string, ending string) (*gzip.Writer, func() error) {
	outName := filepath.Base(inspiration)
	newPath := filepath.Join(outPath, outName+"."+ending)

	out, err := os.Create(newPath)
	if err != nil {
		log.Fatal(err)
	}
	writer := gzip.NewWriter(out)
	return writer, func() error { writer.Flush(); out.Close(); return nil }

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
			log.Fatalf("Appeared to be a directory but had trouble: %s", dirErr)
		}
		// we need to loop over the files
	} else {
		filesToCheck = append(filesToCheck, inFile)
	}
	return filesToCheck
}
