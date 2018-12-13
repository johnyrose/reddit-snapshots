package reddit_snapshot_catcher

import (
	"github.com/jzelinskie/reddit"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type SubredditSnapshot struct {
	Subreddit string
	Time      time.Time
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

func TakeSnapshot(client RedditClient, subreddit string, sort geddit.PopularitySort) SubredditSnapshot {
	s := client.generateSession()
	opt := geddit.ListingOptions{}
	items, err := s.SubredditSubmissions(subreddit, sort, opt)
	if err != nil {
		log.Fatal(err)
	}
	snapshot := SubredditSnapshot{
		Subreddit: subreddit,
		Time:      time.Now(),
		Posts:     items,
	}
	return snapshot
}
