# Profiling what went wrong

I was really surprised last night that protobufs took longer to read than JSON but after profiling, I think I have a few reasons:

- First, protobuf takes more time to deserialize than I thought it would.  That's odd.  Maybe the golang library isn't as efficient as it could be or maybe the JSON library is just super optimized.
- Second, I am reading both files from gzip.  Maybe I am putting too much pressure on the CPU by having it decompress and deserialize - but is it possible that this is the same for JSON?
- I notice in the profiles (which I include in ./profiles/ directory) that there are just a number more steps for protobuf and that the code path seems to be far more complex.  This could be the reason as well.

I think I'll run one more test - that is, unzipping the files and running the comparative analysis to see what's the fastest.  

At this point though, the lesson is that "good ideas" can sometimes cost a lot of time.  I'm sure I will use this information in the future, but like everything, it came at a cost.

```
~/PersCode/reddit_stats:jamesfe$ make compareunzipped
time ./reddit_stats --filename ~/PersCode/reddit_stats/compare/1m_sample_data.json --cv 100000 --maxlines 1000000 --informat json
2017/06/06 09:16:09 Entering analysis loop.
2017/06/06 09:16:09 Read 0 lines
2017/06/06 09:16:09 Read 100000 lines
2017/06/06 09:16:10 Read 200000 lines
2017/06/06 09:16:10 Read 300000 lines
2017/06/06 09:16:10 Read 400000 lines
2017/06/06 09:16:10 Read 500000 lines
2017/06/06 09:16:10 Read 600000 lines
2017/06/06 09:16:10 Read 700000 lines
2017/06/06 09:16:10 Read 800000 lines
2017/06/06 09:16:10 Read 900000 lines
2017/06/06 09:16:10 Max lines reached
2017/06/06 09:16:10 Output written to ./output/output_1496733370.json

real    0m0.953s
user    0m0.809s
sys     0m0.194s
time ./reddit_stats --filename ~/PersCode/reddit_stats/compare/1m_sample_data.json.protodata --cv 100000 --maxlines 1000000 --informat proto
2017/06/06 09:16:10 Delimiter set
2017/06/06 09:16:10 Entering analysis loop.
2017/06/06 09:16:10 Read 0 lines
2017/06/06 09:16:11 Read 100000 lines
2017/06/06 09:16:11 Read 200000 lines
2017/06/06 09:16:11 Read 300000 lines
2017/06/06 09:16:12 Read 400000 lines
2017/06/06 09:16:12 Read 500000 lines
2017/06/06 09:16:12 Read 600000 lines
2017/06/06 09:16:12 Read 700000 lines
2017/06/06 09:16:13 Read 800000 lines
2017/06/06 09:16:13 Read 900000 lines
2017/06/06 09:16:13 Max lines reached
2017/06/06 09:16:13 Output written to ./output/output_1496733373.json

real    0m3.111s
user    0m3.162s
sys     0m0.176s
```

Here is an analysis of comparing unziped files.  We can see here that reading unzipped JSON takes 1/3 as long as reading the unzipped protobuf.  I'm going to go out on a limb here and say that for this project, protobuf is not the right way. 

Did I learn a lot?  For sure.  But sadly I burned about 4 days of free time on this and I have come to the conclusion that I won't use protobufs.
