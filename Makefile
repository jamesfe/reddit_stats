.PHONY: build
build:
	go build -o reddit_stats

.PHONY: smalltest
smalltest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03 --cv 10 --maxlines 100

.PHONY: medtest
medtest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03 --cv 10000 --maxlines 100000

.PHONY: bigtest
bigtest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03 --cv 1000000 --maxlines 10000000

.PHONY: dirtest
dirtest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/dir_test/ -cv 100000 --maxlines 0