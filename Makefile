.PHONY: build
build:
	go build -o reddit_stats ./src/cmd/analyze
	go build -o reddit_filter ./src/cmd/filter
	go build -o full_analyze ./src/cmd/top
	go build -o profile_users ./src/cmd/profiles

.PHONY: buildt
buildt:
	make build
	make test

.PHONY: test
test:
	go test ./src/cmd/analyze -coverprofile=./coverage/cmd_analyze.out
	go test ./src/cmd/top -coverprofile=./coverage/cmd_top.out
	go test ./src/cmd/filter -coverprofile=./coverage/cmd_filter.out
	go test ./src/cmd/profiles -coverprofile=./coverage/cmd_profiles.out
	go test ./src/analysis/ -coverprofile=./coverage/analysis.out
	go test ./src/utils/ -coverprofile=./coverage/utils.out


ifeq (showcoverage,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif

# Show the coverage for a ceratin output file
.PHONY: showcoverage
showcoverage:
	go tool cover -html=./coverage/$(RUN_ARGS).out

.PHONY: tinytestmulti
tinytestmulti:
	time ./full_analyze --config ./configs/tiny.json

.PHONY: tinytest
tinytest:
	time ./reddit_stats --config ./configs/tiny.json

.PHONY: medtestmulti
medtestmulti:
	time ./full_analyze --config ./configs/medium.json

.PHONY: medtest
medtest:
	time ./reddit_stats --config ./configs/medium.json

.PHONY: bigtestmulti
bigtestmulti:
	time ./full_analyze --config ./configs/big.json

.PHONY: bigtest
bigtest:
	time ./reddit_stats --config ./configs/big.json

.PHONY: smfilter
smfilter:
	time ./reddit_filter --filename ~/PersCode/reddit_donald/data/RC_2017-03.gz --cv 1000 --maxlines 10000 --output ./filters/

.PHONY: filterall
filterall:
	time ./reddit_filter --filename ~/PersCode/reddit_donald/data/ --cv 1000000 --maxlines 10000000000 --output ./filters/

.PHONY: analyze
analyze:
	time ./reddit_stats --config ./configs/complete.json

.PHONY: analyzemulti
analyzemulti:
	time ./full_analyze --config ./configs/complete.json
