package storer

import (
	"github.com/Ripolak/reddit-snapshots/catcher"
)

func StoreItem(subreddit catcher.SubredditSnapshot, mongoUrl string, dbName string, collection string) {
	c := MongoClient{
		Url: mongoUrl,
	}
	data := subreddit.ToBsonM()
	c.insertToDatabase(dbName, collection, data)
}
