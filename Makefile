.PHONY: build
build:
	make proto
	go build -o reddit_stats
	go build -o convert ./src/cmd/convert

.PHONY: rbuild
rbuild:
	make proto
	go build -o reddit_stats -race

.PHONY: tinytest
tinytest:
	time ./convert --input ~/PersCode/reddit_donald/data/RC_2016-04.gz --outdir ./protoout/ --from json --numlines 5
	# time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz --cv 1 --maxlines 3 --purpose simple --informat json
	time ./reddit_stats --filename ./protoout/RC_2016-04.gz.protodata --cv 1 --maxlines 5 --informat proto

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


.PHONY: proto
proto:
	protoc --go_out=./ ./reddit_proto/*.proto
