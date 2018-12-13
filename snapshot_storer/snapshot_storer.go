package snapshot_storer

import "reddit-take/reddit_snapshot_taker"

func StoreItem(subreddit reddit_snapshot_taker.SubredditSnapshot, mongoUrl string, dbName string, collection string) {
	c := MongoClient{
		Url: mongoUrl,
	}
	data := subreddit.ToBsonM()
	c.insertToDatabase(dbName, collection, data)
}
