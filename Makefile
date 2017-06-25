.PHONY: build
build:
	go build -o reddit_stats ./src/cmd/analyze
	go build -o reddit_filter ./src/cmd/filter

.PHONY: smalltest
smalltest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz --cv 10 --maxlines 100


.PHONY: medtest
medtest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz --cv 10000 --maxlines 100000


.PHONY: bigtest
bigtest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz --cv 1000000 --maxlines 10000000


.PHONY: dirtest
dirtest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/dir_test/ -cv 100 --maxlines 1000


.PHONY: smfilter
smfilter:
	time ./reddit_filter --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz --cv 1000 --maxlines 10000 --output ./filters/

.PHONY: filterall
filterall:
	time ./reddit_filter --filename ~/PersCode/reddit_donald/data/ --cv 1000000 --maxlines 10000000000 --output ./filters/
