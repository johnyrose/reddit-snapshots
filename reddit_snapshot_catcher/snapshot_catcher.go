package reddit_snapshot_catcher

import (
	"github.com/jzelinskie/reddit"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type SubredditSnapshot struct {
	Subreddit string
	Time      string
	Posts     []*geddit.Submission
}

func (s SubredditSnapshot) ToBsonM() bson.M {
	data, err := bson.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	m := make(bson.M)
	err = bson.Unmarshal(data, &m)
	if err != nil {
		log.Fatal(err)
	}
	return m
}
