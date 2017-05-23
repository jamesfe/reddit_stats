
## Protobuf Conversions Continue

Today I worked on serializing things into protobuf, which has held its own surprises.  Here are the current file sizes for my data, as you can see it's pretty big: 

```
~/PersCode/gosource/src/github.com/jamesfe/reddit_stats:jamesfe$ ls -l ~/PersCode/reddit_donald/data/
total 467890024
-rw-r--r--  1 jamesfe  staff   8986019237 Sep  9  2016 RC_2016-01.gz
-rw-r--r--  1 jamesfe  staff   8745465885 Feb  1 17:43 RC_2016-02.gz
-rw-r--r--  1 jamesfe  staff   9276314460 Feb  1 19:21 RC_2016-03.gz
-rw-r--r--  1 jamesfe  staff   9190059692 Sep  9  2016 RC_2016-04.gz
-rw-r--r--  1 jamesfe  staff   9401909195 Sep  9  2016 RC_2016-05.gz
-rw-r--r--  1 jamesfe  staff   9539355033 Sep  9  2016 RC_2016-06.gz
-rw-r--r--  1 jamesfe  staff   9910922464 Sep  9  2016 RC_2016-07.gz
-rw-r--r--  1 jamesfe  staff  10261082245 Sep 12  2016 RC_2016-08.gz
-rw-r--r--  1 jamesfe  staff   9558986330 Oct 13  2016 RC_2016-09.gz
-rw-r--r--  1 jamesfe  staff   9872557507 Feb  1 19:22 RC_2016-10.gz
-rw-r--r--  1 jamesfe  staff   9903567748 Feb  1 17:47 RC_2016-11.gz
-rw-r--r--  1 jamesfe  staff  10180125098 Jan 12 18:21 RC_2016-12.gz
-rw-r--r--  1 jamesfe  staff  11088597575 Feb  9 06:59 RC_2017-01.gz
-rw-r--r--  1 jamesfe  staff   9967821001 Mar 13 03:15 RC_2017-02.gz
-rw-r--r--  1 jamesfe  staff  42376471592 Apr 14 07:23 RC_2017-03
-rw-r--r--  1 jamesfe  staff   7907014107 Apr 14 07:23 RC_2017-03.bz2
-rw-r--r--  1 jamesfe  staff  11234784813 Apr 14 07:23 RC_2017-03.gz
-rw-r--r--  1 jamesfe  staff  42158585783 May 12 01:07 RC_2017-04
-rwxr-xr-x  1 jamesfe  staff          930 May 17 20:45 get.sh
-rwxr-xr-x  1 jamesfe  staff          717 Feb  1 22:18 zipconv.sh
```
