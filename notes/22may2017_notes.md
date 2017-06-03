After the disaster of concurrency, which in theory sounds fine but would take too long for me to develop, I removed that code.  Everything has been refactored and I've separated out some concerns so I can migrate everything to Protobuf which should show some good speed improvements.

Also, analytically I have removed the "[deleted]" author from my analysis since there is too much ambiguity there, although I may use this as a measure of how much moderation or comment deleting (obvious) is happening.

Here are the results of a bigtest:

```
~/PersCode/gosource/src/github.com/jamesfe/reddit_stats:jamesfe$ make bigtest
time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03 --cv 1000000 --maxlines 10000000
2017/05/22 08:59:43 reading /Users/jamesfe/PersCode/reddit_donald/data/RC_2017-03
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

# Protobuf conversion
I'm going to convert all the data to protobuf for a few reasons:
1. Disk space - sure we compress the data and that saves us a bit on disk space, but we are still reading and decoding the same bytes over and over.  I wonder if protobuf will help with this.
2. Speed: We aren't going to read as many bytes from the disk
3. Move more work from the disk to the processor: I believe the disk (specifically, read) is taking on too much work
4. We can make a single pass over the data and encode it once, throwing away data we don't need.
5. We are going to prep the data as well by creating a subset of the data in separate files.  This way we can scan the target subreddit directly while still having access to the greater body of data.

There are drawbacks:
1. We may lose some data if I throw it away for error or marshalling reasons during the conversion.  More risk here too since I am going to throw away the gzip files after I am confident in the conversion.
2. We will possibly be CPU bound after this (mitigating factor: the disaster of goroutines returns?)
3. We have a hard and fast schema and if new data becomes available (or if fields go away) in the future we will have to flex.
