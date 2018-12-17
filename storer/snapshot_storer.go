package storer

import (
	"github.com/Ripolak/reddit-snapshots/catcher"
)

type SnapshotStorer interface {
	StoreItem(subreddit catcher.SubredditSnapshot)
}

type DatabaseStorer struct {
	MongoClient mongoClient
	DbName      string
	Collection  string
}

func (d DatabaseStorer) StoreItem(subreddit catcher.SubredditSnapshot) {
	data := subreddit.ToBsonM()
	d.MongoClient.insertToDatabase(d.DbName, d.Collection, data)
}
