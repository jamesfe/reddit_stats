# About Reddit Stats

## General Purpose

This repo is a collection of Golang programs that you can use to analyze reddit comments.  The outputscan be controlled by a variety of configuration options: here are a few analysis pieces I've done so far: 

- List the number of unique commenters 
- List the number of deleted comments

These can be aggregated by day, week, or any other date aggregation function.
Also, these can target a single subreddit (faster) or they can target lists of subreddits (much slowerr).

## Other Purposes

I have written Go before but not as much as I wanted and I've never set up a larger project with modules like this one.  (This, for the record, is not a "big" project, but larger than fun scripts I had written before.)  So this serves a few purposes:

- Learn about Go, build tools, best practices
- Practice learning about Go because reading is not alone enough
- Learn about optimizing and profiling Go programs (./notes/profiles/)
- Learn about parsing and aggregating data over large datasets (1.3bn+ lines of JSON in play here)

## Non-Purposes & Alternatives to reddit_stats

- Build a replacement for anything
- Do anything faster than an RDBMS
- If you are an enterprise user you should probably just use Hive or BigQuery

## Notes

I keep a journal of thoughts that occur to me while coding here: (./notes/)

## Contact Me

I love talking to other people about all sorts of things.  Feel free to reach out!
Twitter: [@jimmysthoughts](https://twitter.com/jimmysthoughts)
