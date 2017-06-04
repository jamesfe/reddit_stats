# Notes 2 Jun 2017

I was sick and/or on vacation for about a week but the results are in for converting one large file to protobuf.

```
-rw-r--r--   1 jamesfe  staff   8986019237 Sep  9  2016 RC_2016-01.gz // json
-rw-r--r--  1 jamesfe  staff  7733457293 May 24 10:06 RC_2016-01.gz // protobuf
```

We see that the protobuf file is 86.06% as large as the before JSON file.  This is good - about as good as we expected it would be.  We expected a 14.93% savings and that is almost what we got (13.93%).

Now we need to run a test to see if reading the files is just as good.  Let's run the same author analysis on both files.


## To Do
- Write protobuf reading routines
- Add a flag to read one or the other
- Run the regular simple author analysis on both sides: proto and non-proto
