package reddit_proto

import (
	"encoding/json"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("reddit_proto")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}`,
)

func GetIntTimestamp(v interface{}) int64 {
	/* Sometimes the timestamps we get are float64, null, or strings. Here we check. */
	var retVal int64 = 0
	switch v.(type) {
	case float64:
		retVal = int64(v.(float64))
	case string:
		parsed, err := strconv.ParseInt(v.(string), 10, 32)
		if err != nil {
			retVal = 0
		} else {
			retVal = int64(parsed)
		}
	default:
		retVal = 0
	}
	return retVal
}

func ConvertLineToProto(line []byte, outdata *[]byte) bool {
	var rawJsonMap interface{}
	jumerr := json.Unmarshal(line, &rawJsonMap)
	var newComment Comment
	v := rawJsonMap.(map[string]interface{})

	if jumerr != nil {
		return false
	}

	for k, v := range v {
		switch vv := v.(type) {
		case string:
			switch k {
			case "author":
				if vv == "" {
					return false // problem
				}
				newComment.Author = &vv
			case "body":
				newComment.Body = &vv
			case "author_flair_css_class":
				newComment.AuthorFlairCssClass = &vv
			case "author_flair_text":
				newComment.AuthorFlairText = &vv
			case "edited":
				newComment.Edited = &vv
			case "created_utc":
				log.Errorf("Whoa! A string timestamp!")
				return false
			case "link_id":
				newComment.LinkId = &vv
			case "parent_id":
				newComment.ParentId = &vv
			case "subreddit":
				newComment.Subreddit = &vv
			case "subreddit_id":
				newComment.SubredditId = &vv
			}
		case float64: // all the numbers are by default float64 from json
			item64 := int64(vv)
			item32 := int32(vv)
			switch k {
			case "created_utc":
				if item64 == 0 {
					return false
				}
				newComment.CreatedUTC = &item64
			case "controversiality":
				newComment.Controversiality = &item32
			case "gilded":
				newComment.Gilded = &item32
			case "retrieved_on":
				newComment.RetrievedOn = &item64
			case "ups":
				newComment.Ups = &item32
			case "score":
				newComment.Score = &item32

			}
		}
	}
	outd, err := proto.Marshal(&newComment)
	outdata = &outd
	if err != nil {
		return false
	}
	return true
}
