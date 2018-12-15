package main

import (
	"github.com/Ripolak/reddit-snapshots/reddit_snapshot_catcher"
	"github.com/Ripolak/reddit-snapshots/snapshot_storer"
	"os"
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
	reddit              = reddit_snapshot_catcher.RedditClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Username:     username,
		Password:     password,
	}
)

func main() {
	snapshotConfig := snapshot_storer.LoadConfiguration(dbUrl, dbName, configCollection)
	subreddits := snapshotConfig.Subreddits
	fetchSnapshots(subreddits)
}
