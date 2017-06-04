# Multiple Executables

I realized that the big analysis file I was working with was too much and have begun to split things out into separate modules.  First, I want a module that just converts to and from proto - I should have done this earlier. Now it's done.

In the refactor, maybe even in the writing, I was calling `flush()` about a zillion times too many.  This might have thrown me off (did it write a bunch of gzip junk to the file?) but I will check later.

For now though I think I am ready to convert to proto and unmarshal successfully.

```
~/PersCode/reddit_stats:jamesfe$ make fullconvert
time ./convert --input ~/PersCode/reddit_donald/data/RC_2016-01.gz --outdir ./protoout/ --from json --numlines 1000000000

real    43m51.828s
user    43m35.213s
sys     1m23.571s
~/PersCode/reddit_stats:jamesfe$ ls -l ./protoout/
-rw-r--r--  1 jamesfe  staff  7738771935 Jun  3 21:39 RC_2016-01.gz.protodata.gz
```

So I think this is our file.  We flushed a lot fewer times but it is still 5mb larger.  Alas.

Let's compare the two types of analysis in terms of time though.
