package protoanalysis

import (
	"strings"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/jamesfe/reddit_stats/data_types"
	"github.com/jamesfe/reddit_stats/reddit_proto"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("proto-analysis")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
)

func ProtoSingleLineAnalysis(line []byte, result *data_types.AuthorDateTuple) bool {
	/* Take bytes, convert to proto object, check if author & date are valid and return object if valid. */

	comment := &reddit_proto.Comment{}
	err := proto.Unmarshal(line, comment)
	if err != nil {
		log.Errorf("Unmarshall error", err)
	} else {
		log.Infof("Proto: %#v", comment)
		if strings.ToLower(*comment.Subreddit) == "the_donald" {
			if *comment.CreatedUTC != 0 && *comment.Author != "[deleted]" {
				result.AuthorName = *comment.Author
				result.AuthorDate = time.Unix(int64(*comment.CreatedUTC), 0).Format("02-01-2006")
				return true
			}
		}
	}
	return false
}
