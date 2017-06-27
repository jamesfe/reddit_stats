# 27 June 2017 Notes

## Date Aggregation

I'm going to start bucketing events by weeks to make the graph a little more understandable.  To do this I've refactored out the date aggregation function and it's now a parameter of the analysis function.  I've added one that returns WW-YYYY format strings for more wide aggregation.  

## Refactoring

A new month has come and I need to provide my usual graph of unique authors, but unfortunately I have made more commits on top of that.  So I need to go back and refactor the aggregation function so it'll provide those stats again.

## Next Steps

I need to being profiling users of target reddits so I can build a model for them and begin to measure their level of interaction with the greater reddit community - in effect, to find if they're isolated.  Coming soon.

## News

I published an article about this process here: [https://medium.com/@jimmysthoughts/aggregating-reddit-comments-34ab44e48cb1]
