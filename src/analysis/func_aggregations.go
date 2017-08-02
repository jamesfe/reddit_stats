package analysis

import (
	"github.com/jamesfe/reddit_stats/src/data_types"
	"github.com/jamesfe/reddit_stats/src/utils"
)

func AggregateUniqueAuthors(analysisResults map[string]map[string]int) UniqueAuthorsPerDayResult {
	resultMap := NewUniqueAuthorsPerDayResult()
	for k, v := range analysisResults {
		(*resultMap).AuthorsPerDay[k] = len(v)
	}
	return *resultMap
}

func CreateActiveUserMap(input map[string]*data_types.UserLongevityResult, start int, end int, delta int, minSecondsDiff int, dateFunc data_types.DateToString) map[string]int {
	/*
	   First we create a map with all the dates we want.  We need start, end, and delta and aggregate function.
	   Second we iterate over the list of items.  We generate a date, round it to delta and then add delta until it's larger than end.
	*/
	rv := make(map[string]int)
	// Make the original map
	for a := start; a < end; a += delta {
		rv[dateFunc(a)] = 0
	}

	// For each item, add a bunch of things.
	for _, element := range input {
		if (element.LastPost - element.FirstPost) >= minSecondsDiff {
			for b := element.FirstPost - (element.FirstPost % delta); b < element.LastPost; b += delta {
				rv[dateFunc(b)] += 1
			}
		}
	}
	return rv
}

func AggregateByAuthorLongevity(input map[string]*data_types.UserLongevityResult, minSecondsDiff int) []*data_types.TimePeriod {
	var rv []*data_types.TimePeriod
	for _, element := range input {
		if (element.LastPost - element.FirstPost) >= minSecondsDiff {
			newObject := &data_types.TimePeriod{
				StartDate: utils.GetDayString(element.FirstPost),
				EndDate:   utils.GetDayString(element.LastPost)}
			rv = append(rv, newObject)
		}
	}
	return rv
}

func AggregateByDeletedCommentCounts(analysisResults map[string]map[string]int) map[string]data_types.DeletedTuple {
	/* Count up the number of deleted, total, and not-deleted comments per time period and return them in a map. */
	var today_sum int
	var deleted_sum int
	var total_sum int
	outputMap := make(map[string]data_types.DeletedTuple)
	for key, element := range analysisResults {
		today_sum = 0
		deleted_sum = 0
		for author, count := range element {
			if author != "[deleted]" {
				today_sum += count
			} else {
				deleted_sum = count
			}
		}
		total_sum = today_sum + deleted_sum
		d := &data_types.DeletedTuple{TodayTotal: today_sum, Deleted: deleted_sum, Total: total_sum}
		outputMap[key] = *d
	}
	return outputMap
}
