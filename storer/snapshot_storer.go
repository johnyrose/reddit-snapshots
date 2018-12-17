package storer

import (
	"github.com/Ripolak/reddit-snapshots/catcher"
)

type SnapshotStorer interface {
	StoreItem(subreddit catcher.SubredditSnapshot)
}

type DatabaseStorer struct {
	MongoClient MongoClient
	DbName      string
	Collection  string
}

func (d DatabaseStorer) StoreItem(subreddit catcher.SubredditSnapshot) {
	data := subreddit.ToBsonM()
	d.MongoClient.insertToDatabase(d.DbName, d.Collection, data)
}

func StoreItem(subreddit catcher.SubredditSnapshot, mongoUrl string, dbName string, collection string) {
	c := MongoClient{
		Url: mongoUrl,
	}
	data := subreddit.ToBsonM()
	c.insertToDatabase(dbName, collection, data)
}
