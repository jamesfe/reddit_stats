# More protobuf

I ran the program yesterday: it took about 30 minutes to convert everything from JSON to Protobuf but my disk filled up.

Also I was surprised - the protobuf was larger than the gzipped JSON.  So, let's GZIP the protobuf.  We will see.

In general though, I am seeing an overall savings of 51% space and an average of 46% for each JSON line.  This was on a test of 1,000 lines of JSON.

Here are the results for zipping everything:

Inputs (before, after zipping) and output:
```
505051 1k_sample_data.json
134552 1k_sample_data.json.gz
114458 1k_sample_data.protodata
```

This is 1,000 JSON records, so let's run some numbers:
1. An overall savings from unzipped data of 77.33%
2. An overall savings from zipped data of 14.93% (this is what we will realize)

So here is a table of file sizes:

| Bytes Before | File | Estimated Bytes After | Estimated MBytes Saved |
| --- | --- | --- | --- |
| 8986019237 | RC_2016-01.gz | 7644046835 | 1341.97 |
| 8745465885 | RC_2016-02.gz | 7439417728 | 1306.05 |
| 9276314460 | RC_2016-03.gz | 7890989360 | 1385.33 |
| 9190059692 | RC_2016-04.gz | 7817615882 | 1372.44 |
| 9401909195 | RC_2016-05.gz | 7997827773 | 1404.08 |
| 9539355033 | RC_2016-06.gz | 8114747446 | 1424.61 |
| 9910922464 | RC_2016-07.gz | 8430824985 | 1480.1 |
| 10261082245 | RC_2016-08.gz | 8728691893 | 1532.39 |
| 9558986330 | RC_2016-09.gz | 8131447004 | 1427.54 |
| 9872557507 | RC_2016-10.gz | 8398189451 | 1474.37 |
| 9903567748 | RC_2016-11.gz | 8424568622 | 1479 |
| 10180125098 | RC_2016-12.gz | 8659824889 | 1520.3 |
| 11088597575 | RC_2017-01.gz | 9432626057 | 1655.97 |
| 9967821001 | RC_2017-02.gz | 8479226292 | 1488.59 |
| 11234784813 | RC_2017-03.gz | 9556981688 | 1677.8 |
| 10879763818 | RC_2017-04.gz | 9254979540 | 1624.78 |

So the next step is to see if saving space has any correlation with being more performant.  I'm going to run a medium sized analysis on the JSON data and record the results, then I will run the same analysis on the Protobuf data and we will see.

Here are the results (truncated a bit)
```
~/PersCode/reddit_stats:jamesfe$ make medtest
time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz --cv 10000 --maxlines 100000 --purpose simple
2017/05/24 09:14:55 Entering analysis stream.
2017/05/24 09:14:55 Read 0 lines
2017/05/24 09:14:55 Read 90000 lines
2017/05/24 09:14:55 Max lines reached
2017/05/24 09:14:55 Output written to ./output/output_1495610095.json

real    0m0.591s
user    0m0.547s
sys     0m0.031s
~/PersCode/reddit_stats:jamesfe$ make bigtest
time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz --cv 1000000 --maxlines 10000000 --purpose simple
2017/05/24 09:15:00 Entering analysis stream.
2017/05/24 09:15:00 Read 0 lines
2017/05/24 09:15:52 Read 9000000 lines
2017/05/24 09:15:57 Max lines reached
2017/05/24 09:15:57 Output written to ./output/output_1495610157.json

real    0m57.031s
user    0m58.167s
sys     0m1.450s
```

And the output directory:
```
-rw-r--r--   1 jferrara  staff   16828 May 24 09:14 output_1495610095.json // medium
-rw-r--r--   1 jferrara  staff  670500 May 24 09:15 output_1495610157.json // big
```
