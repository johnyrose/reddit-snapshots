package manager

import (
	"github.com/jzelinskie/reddit"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/mgo.v2/bson"
	"log"
	"sync"

	"github.com/Ripolak/reddit-snapshots/catcher"
	"github.com/Ripolak/reddit-snapshots/storer"
)

type dbConfig struct {
	DbUrl               string
	DbName              string
	SnapshotsCollection string
	ConfigCollection    string
}

type redditConfig struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

func Entrypoint() {

	var db dbConfig
	err := envconfig.Process("", &db)
	if err != nil {
		log.Fatal(err)
	}
	var reddit redditConfig
	err = envconfig.Process("", &reddit)
	if err != nil {
		log.Fatal(err)
	}

	snapshotConfig := storer.LoadConfiguration(dbUrl, dbName, configCollection)
	subreddits := snapshotConfig.Subreddits
	fetchSnapshots(subreddits)
}

func fetchSnapshots(subreddits []bson.M) {
	var wg sync.WaitGroup
	wg.Add(len(subreddits))
	ch := make(chan catcher.SubredditSnapshot, len(subreddits))
	for _, subreddit := range subreddits {
		go takeSnapshot(&wg, subreddit["subreddit"].(string), geddit.PopularitySort(subreddit["sort"].(string)), ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	storeSnapshots(ch)
}

func takeSnapshot(wg *sync.WaitGroup, subreddit string, sort geddit.PopularitySort, ch chan catcher.SubredditSnapshot) {
	defer wg.Done()
	snapshot := catcher.TakeSnapshot(reddit, subreddit, sort)
	ch <- snapshot
}

func storeSnapshots(ch chan catcher.SubredditSnapshot) {
	for msg := range ch {
		storer.StoreItem(msg, dbUrl, dbName, snapshotsCollection)
	}
}
