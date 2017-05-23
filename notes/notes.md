# Undefined notes from random days:

just reading 10k without scanning the string

2017/04/17 19:56:22 reading ./data/10k_sample_data.json.gz

real    0m0.214s
user    0m0.210s
sys     0m0.012s

read the string first, lower, check for donald:

2017/04/17 19:54:59 reading ./data/10k_sample_data.json.gz

real    0m0.158s
user    0m0.148s
sys     0m0.013s

all comments, no string scanning for The_donald
~/PersCode/gosource/src/github.com/jamesfe/reddit_stats:jamesfe$ time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz
2017/04/17 21:12:20 reading /Users/jamesfe/PersCode/reddit_donald/data/RC_2017-03.gz
2017/04/17 21:13:31 Fatal Scanning Error: %s
 <nil>

 real    1m11.000s
 user    1m14.769s
 sys     0m1.914s

with scanning
 ~/PersCode/gosource/src/github.com/jamesfe/reddit_stats:jamesfe$ go build
 ~/PersCode/gosource/src/github.com/jamesfe/reddit_stats:jamesfe$ time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz
 2017/04/17 21:14:22 reading /Users/jamesfe/PersCode/reddit_donald/data/RC_2017-03.gz
 2017/04/17 21:15:14 Fatal Scanning Error <nil>

 real    0m51.466s
 user    0m52.957s
 sys     0m1.154s


scanning and unmarshalling the body
~/PersCode/gosource/src/github.com/jamesfe/reddit_stats:jamesfe$ time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz
2017/04/17 21:22:22 reading /Users/jamesfe/PersCode/reddit_donald/data/RC_2017-03.gz
2017/04/17 21:23:16 Fatal Scanning Error <nil>

real    0m53.629s
user    0m55.094s
sys     0m1.249s


with a counter every step and match: 
~/PersCode/gosource/src/github.com/jamesfe/reddit_stats:jamesfe$ time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz
2017/04/17 21:26:14 reading /Users/jamesfe/PersCode/reddit_donald/data/RC_2017-03.gz
81407 lines matched out of 3423896
real    0m53.361s
user    0m54.806s
sys     0m1.245s
~/PersCode/gosource/src/github.com/jamesfe/reddit_stats:jamesfe$ time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz
2017/04/17 21:28:11 reading /Users/jamesfe/PersCode/reddit_donald/data/RC_2017-03.gz
81407 lines matched out of 3423896
real    0m54.128s
user    0m55.066s
sys     0m1.444s


^C~/PersCode/gosource/src/github.com/jamesfe/reddit_stats:jamesfe$ ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz -cv 100000
2017/04/18 21:44:45 reading /Users/jamesfe/PersCode/reddit_donald/data/RC_2017-03.gz
2017/04/18 21:44:48 Read 100000 lines
2017/04/18 21:44:50 Read 200000 lines
2017/04/18 21:44:52 Read 300000 lines
2017/04/18 21:44:55 Read 400000 lines
2017/04/18 21:44:57 Read 500000 lines
2017/04/18 21:44:59 Read 600000 lines
2017/04/18 21:45:02 Read 700000 lines
2017/04/18 21:45:04 Read 800000 lines
2017/04/18 21:45:06 Read 900000 lines
2017/04/18 21:45:09 Read 1000000 lines
2017/04/18 21:45:11 Read 1100000 lines
2017/04/18 21:45:14 Read 1200000 lines
2017/04/18 21:45:17 Read 1300000 lines
2017/04/18 22:16:01 Read 79300000 lines
2017/04/18 22:16:03 Read 79400000 lines
2017/04/18 22:16:05 Read 79500000 lines
2017/04/18 22:16:08 Read 79600000 lines
2017/04/18 22:16:10 Read 79700000 lines
1822765, 1362541 (initial, final) lines matched out of 797231062017/04/18 22:16:11 READLINE:  EOF
~/PersCode/gosource/src/github.com/jamesfe/reddit_stats:jamesfe$ python
Python 2.7.12 (default, Sep 12 2016, 09:26:13)
[GCC 4.2.1 Compatible Apple LLVM 7.3.0 (clang-703.0.31)] on darwin
Type "help", "copyright", "credits" or "license" for more information.
>>> 1822765 - 1362541
460224










19 may 2017
today I changed a lot of stuff to go routines.
I changed the reading from files to a goroutine based thingy and that's simpler in terms of coding
A little bit slower and more bugs with concurrency if you decide to aggregate things as you read/process them
Probably smarter to use protobuf

Major concurrency issues when reading & aggregating simultaneously

Here are the results of a simple author count (before is with goroutines + race protection code and after
is with just handling the goroutines)

~/PersCode/gosource/src/github.com/jamesfe/reddit_stats:jamesfe$ make bigtest
time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03 --cv 1000000 --maxlines 10000000
2017/05/19 22:47:06 reading /Users/jamesfe/PersCode/reddit_donald/data/RC_2017-03
2017/05/19 22:47:06 Entering analysis stream.
2017/05/19 22:47:06 Read 0 lines
2017/05/19 22:47:17 Read 1000000 lines
2017/05/19 22:47:28 Read 2000000 lines
2017/05/19 22:47:41 Read 3000000 lines
2017/05/19 22:48:07 Read 4000000 lines
2017/05/19 22:48:42 Read 5000000 lines
2017/05/19 22:48:48 Read 6000000 lines
2017/05/19 22:49:38 Read 7000000 lines
2017/05/19 22:49:44 Read 8000000 lines
2017/05/19 22:49:49 Read 9000000 lines
2017/05/19 22:51:09 Output written to ./output/output_1495227069.json

real    4m5.655s
user    4m27.552s
sys     2m49.134s

~/PersCode/gosource/src/github.com/jamesfe/reddit_stats:jamesfe$ make bigtest
time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03 --cv 1000000 --maxlines 10000000
2017/05/19 22:52:58 reading /Users/jamesfe/PersCode/reddit_donald/data/RC_2017-03
2017/05/19 22:52:58 Entering analysis stream.
2017/05/19 22:52:58 Read 0 lines
2017/05/19 22:53:01 Read 1000000 lines
2017/05/19 22:53:04 Read 2000000 lines
2017/05/19 22:53:07 Read 3000000 lines
2017/05/19 22:53:10 Read 4000000 lines
2017/05/19 22:53:13 Read 5000000 lines
2017/05/19 22:53:17 Read 6000000 lines
2017/05/19 22:53:20 Read 7000000 lines
2017/05/19 22:53:23 Read 8000000 lines
2017/05/19 22:53:26 Read 9000000 lines
2017/05/19 22:53:29 Output written to ./output/output_1495227209.json

real    0m30.832s
user    1m14.028s
sys     0m8.182s
