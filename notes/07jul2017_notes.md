# Notes for 07 July 2017

## Today's Work

Today my goal is to collect analytical results for all the subreddits.  I will do this like so:

1. Copy the entire main analysis program since we need significant data structure changes (we could do it without but would re-parse the data n times for `n=# of subreddits`).
2. Modify and refactor the `func_author_analysis.go` file to create some new functionality.

## Side Notes

*Testing* - I should be doing this.  I'll look into a framework for testing a lot of my business logic and start to think of edge cases.  I think this will help as the program gets bigger.

*New Idea* - Last night I thought it would be nice to calculate the average longevity of a user over a time period - in effect if User A joins the subreddit, when is his/her first & last comment to the subreddit?  If it is a 100 day period and they are present for 20 days, they are a 20% user.  I should probably weight this with how *frequently* they post since a two-post user could be 100% longevity, but maybe that's the joy of this statistic?
