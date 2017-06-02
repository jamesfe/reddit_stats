.PHONY: build
build:
	make proto
	go build -o reddit_stats

.PHONY: rbuild
rbuild:
	make proto
	go build -o reddit_stats -race

.PHONY: protocomp
protocomp:
	time ./reddit_stats --filename ./test_data/1k_sample_data.json -cv 100 --maxlines 1000 --purpose proto --output ./protoout/

.PHONY: prototest
prototest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/data/ --cv 100 --maxlines 1000 --purpose proto --output ./protoout/

.PHONY: fileproto
fileproto:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/data/ --cv 100000 --maxlines 10000000 --purpose proto --output ./protoout/

.PHONY: medprototest
medprototest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/data/ --cv 10000 --maxlines 100000 --purpose proto --output ./protoout/


.PHONY: tinytest
tinytest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz --cv 1 --maxlines 10 --purpose simple

.PHONY: smalltest
smalltest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz --cv 10 --maxlines 100 --purpose simple


.PHONY: medtest
medtest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz --cv 10000 --maxlines 100000 --purpose simple


.PHONY: bigtest
bigtest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz --cv 1000000 --maxlines 10000000 --purpose simple


.PHONY: dirtest
dirtest:
	time ./reddit_stats --filename ~/PersCode/reddit_donald/dir_test/ -cv 100 --maxlines 1000 --purpose simple


.PHONY: proto
proto:
	protoc --go_out=./ ./reddit_proto/*.proto
