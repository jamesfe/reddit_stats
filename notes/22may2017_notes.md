After the disaster of concurrency, which in theory sounds fine but would take too long for me to develop, I removed that code.  Everything has been refactored and I've separated out some concerns so I can migrate everything to Protobuf which should show some good speed improvements.

Also, analytically I have removed the "[deleted]" author from my analysis since there is too much ambiguity there, although I may use this as a measure of how much moderation or comment deleting (obvious) is happening.

Here are the results of a bigtest:

```
~/PersCode/gosource/src/github.com/jamesfe/reddit_stats:jferrara$ make bigtest
time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03 --cv 1000000 --maxlines 10000000
2017/05/22 08:59:43 reading /Users/jferrara/PersCode/reddit_donald/data/RC_2017-03
2017/05/22 08:59:43 Entering analysis stream.
2017/05/22 08:59:43 Read 0 lines
2017/05/22 08:59:46 Read 1000000 lines
2017/05/22 08:59:48 Read 2000000 lines
2017/05/22 08:59:50 Read 3000000 lines
2017/05/22 08:59:53 Read 4000000 lines
2017/05/22 08:59:55 Read 5000000 lines
2017/05/22 08:59:58 Read 6000000 lines
2017/05/22 09:00:00 Read 7000000 lines
2017/05/22 09:00:03 Read 8000000 lines
2017/05/22 09:00:05 Read 9000000 lines
2017/05/22 09:00:08 Output written to ./output/output_1495436408.json

real    0m24.717s
user    0m29.429s
sys     0m4.412s
```


We can see that they are comparable to the concurrency stuff - there wasn't much gain there since we are mostly waititng for disk I/O.


