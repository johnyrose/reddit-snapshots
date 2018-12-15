package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"os"
	"sync"

	"github.com/Ripolak/reddit-snapshots/reddit_snapshot_catcher"
	"github.com/Ripolak/reddit-snapshots/snapshot_storer"
)

var (
	clientID            = os.Getenv("CLIENT_ID")
	clientSecret        = os.Getenv("CLIENT_SECRET")
	username            = os.Getenv("USERNAME")
	password            = os.Getenv("PASSWORD")
	dbUrl               = os.Getenv("DB_URL")
	dbName              = os.Getenv("DB_NAME")
	snapshotsCollection = os.Getenv("SNAPSHOTS_COLLECTION")
	configCollection    = os.Getenv("CONFIG_COLLECTION")
)

func main() {
	snapshotConfig := snapshot_storer.LoadConfiguration(dbUrl, dbName, configCollection)
	reddit := reddit_snapshot_catcher.RedditClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Username:     username,
		Password:     password,
	}

	subreddits := snapshotConfig.Subreddits
	fetchSnapshots(subreddits)

	for _, subreddit := range subreddits {
		snapshot := reddit_snapshot_catcher.TakeSnapshot(reddit, subreddit["subreddit"].(string), "hot")
		snapshot_storer.StoreItem(snapshot, dbUrl, dbName, snapshotsCollection)
	}
}

func fetchSnapshots(subreddits []bson.M, redditClient reddit_snapshot_catcher.RedditClient) {
	var wg sync.WaitGroup
	wg.Add(len(subreddits) * 2)
	ch := make(chan reddit_snapshot_catcher.SubredditSnapshot, len(subreddits))
	takeSnapshots(subreddits, &wg, redditClient, ch)
	wg.Wait()
}

func takeSnapshots(subreddits []bson.M, wg *sync.WaitGroup, redditClient reddit_snapshot_catcher.RedditClient,
	ch chan reddit_snapshot_catcher.SubredditSnapshot) {
	for _, subreddit := range subreddits {
		go func(subreddit string) {
			defer wg.Done()
			snapshot := reddit_snapshot_catcher.TakeSnapshot(redditClient, subreddit, "hot")
			ch <- snapshot
		}(subreddit["subreddit"].(string))
	}
}
