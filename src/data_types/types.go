package data_types

// A function that takes a timestamp and returns some string interpretation of it.
type DateToString func(int) string

type AuthorDateTuple struct {
	AuthorDate string
	AuthorName string
	Timestamp  int
}

type AuthorSubTuple struct {
	SubReddit  string
	AuthorName string
}
type AuthorDateSubTuple struct {
	AuthorDate string
	SubReddit  string
	AuthorName string
}

type Comment struct {
	Author string `json:"author"`
	// Author_flair_css_class string `json:"author_flair_css_class "`
	// Author_flair_text      string `json:"author_flair_text"`
	// Body             string `json:"body"`
	Controversiality int `json:"controversiality"`
	// Created_utc      string `json:"created_utc"`
	// distinguished - mosty null
	// Edited       string `json:"edited"`
	Gilded int    `json:"gilded"`
	Id     string `json:"id"`
	// Link_id      string `json:"link_id"`
	// Parent_id    string `json:"parent_id"`
	// Retrieved_on int    `json:"retrieved_on"`
	Score     int    `json:"score"`
	Subreddit string `json:"subreddit"`
	// Subreddit_id string `json:"subreddit_id"`
	// Ups          int    `json:"ups"`
}

type DeletedTuple struct {
	TodayTotal int `json:"total_not_deleted"`
	Deleted    int `json:"total_deleted"`
	Total      int `json:"total"`
}

type DeletedByDate struct {
	Dates map[string]DeletedTuple `json:"dates"`
}

type DeletedByDateAndReddit struct {
	Reddits map[string]DeletedByDate `json:"reddits"`
}

type UserLongevityResult struct {
	// Author	string `json:"author"`
	FirstPost int `json:"first_post_timestamp"`
	LastPost  int `json:"last_post_timestamp"`
	NumPosts  int `json:"num_posts"`
}

type TimePeriod struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

/*
For Later Use
type AuthorProfile struct {
	// MinDate            int            `json:"first_post_timestamp"`
	// MaxDate            int            `json:"last_post_timestamp"`
	CommentCountsBySub map[string]int `json:"comments_per_sub"`
}
*/

type JSONList struct {
	Items []string `json:"items"`
}
