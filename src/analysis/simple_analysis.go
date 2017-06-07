package analysis

func newSimpleAnalysisResult() SimpleAnalysisResult {
	return (SimpleAnalysisResult{TotalMatches: 0, TotalFirstMatches: 0, TotalLinesChecked: 0})
}

type SimpleAnalysisResult struct {
	TotalMatches      int
	TotalFirstMatches int
	TotalLinesChecked int
}

type SimpleAnalysisParameter struct {
	LinesToCheck             int  // Max number of lines per file to check
	CheckLines               bool // even check them?
	LineIntervalNotification int  // How many lines between print statements?
	LogLineNotification      bool // whether or not to print notifications at line vals
	Filename                 string
}

/* Opens the file, reads it, counts some things up and returns a set of results.  */
/*
func SimpleFileAnalysis(parameters SimpleAnalysisParameter) (SimpleAnalysisResult, error) {
	inFileReader, f := getFileReader(parameters.Filename)
	defer f()
	results := newSimpleAnalysisResult()

	for {
		var v Comment
		var stuff, err = inFileReader.ReadBytes('\n')
		if err != nil {
			log.Warningf("%d, %d (initial, final) lines matched out of %d", results.TotalFirstMatches, results.TotalMatches, results.TotalLinesChecked)
			return results, err
		}
		if (parameters.CheckLines) && (results.TotalLinesChecked >= parameters.LinesToCheck) {
			log.Errorf("Max lines of %d exceeded: %d", parameters.LinesToCheck, results.TotalLinesChecked)
			return results, nil
		}
		if isDonaldLite(stuff) {
			results.TotalFirstMatches += 1
			newerr := json.Unmarshal(stuff, &v)

			if newerr == nil && isDonaldCertainly(v) {
				results.TotalMatches += 1
			} else {
				return results, newerr
			}
		}
		if parameters.LogLineNotification && results.TotalLinesChecked%parameters.LineIntervalNotification == 0 {
			log.Debugf("Read %d lines", results.TotalLinesChecked)
		}
		results.TotalLinesChecked++
	}
	return results, nil
}
*/
