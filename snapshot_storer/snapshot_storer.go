package snapshot_storer

import (
	"github.com/Ripolak/reddit-snapshots/reddit_snapshot_catcher"
)

func StoreItem(subreddit reddit_snapshot_catcher.SubredditSnapshot, mongoUrl string, dbName string, collection string) {
	c := MongoClient{
		Url: mongoUrl,
	}
	data := subreddit.ToBsonM()
	c.insertToDatabase(dbName, collection, data)
}
