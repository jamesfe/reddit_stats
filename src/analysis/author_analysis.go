package analysis

import "time"

type UniqueAuthorsPerDayResult struct {
	AuthorsPerDay map[string]int `json:"authors_per_day"`
	StartDate     time.Time      `json:"-"`
	EndDate       time.Time      `json:"-"`
}

func NewUniqueAuthorsPerDayResult() *UniqueAuthorsPerDayResult {
	retVal := new(UniqueAuthorsPerDayResult)
	retVal.AuthorsPerDay = map[string]int{}
	return retVal
}
