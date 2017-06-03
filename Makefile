.PHONY: build
build:
	make proto
	go build -o reddit_stats
	go build -o convert ./src/cmd/convert

.PHONY: rbuild
rbuild:
	make proto
	go build -o reddit_stats -race

# Convert a file, then return it to JSON
.PHONY: tinytest
tinytest:
	./convert --input ~/PersCode/reddit_donald/data/RC_2016-04.gz --outdir ./protoout/ --from json --numlines 5
	./convert --input ./protoout/RC_2016-04.gz.protodata.gz --outdir ./protoout/ --from proto --numlines 5

# Convert an entire file
.PHONY: fullconvert
fullconvert:
	time ./convert --input ~/PersCode/reddit_donald/data/RC_2016-01.gz --outdir ./protoout/ --from json --numlines 1000000000

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
