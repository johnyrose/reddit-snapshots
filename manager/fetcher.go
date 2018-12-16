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
	DbUrl               string `split_words:"true"`
	DbName              string `split_words:"true"`
	SnapshotsCollection string `split_words:"true"`
	ConfigCollection    string `split_words:"true"`
}

type redditConfig struct {
	ClientID     string `split_words:"true"`
	ClientSecret string `split_words:"true"`
	Username     string `split_words:"true"`
	Password     string `split_words:"true"`
}

func Entrypoint() {

	var db dbConfig
	err := envconfig.Process("", &db)
	if err != nil {
		log.Fatal(err)
	}
	var redditConfig redditConfig
	err = envconfig.Process("", &redditConfig)
	if err != nil {
		log.Fatal(err)
	}

	snapshotConfig := storer.LoadConfiguration(dbUrl, dbName, configCollection)
	subreddits := snapshotConfig.Subreddits
	fetchSnapshots(subreddits, redditConfig)
}

func fetchSnapshots(subreddits []bson.M, redditConfig redditConfig) {
	reddit := catcher.RedditClient{
		ClientID:     redditConfig.ClientID,
		ClientSecret: redditConfig.ClientSecret,
		Username:     redditConfig.Username,
		Password:     redditConfig.Password,
	}
	var wg sync.WaitGroup
	wg.Add(len(subreddits))
	ch := make(chan catcher.SubredditSnapshot, len(subreddits))
	for _, subreddit := range subreddits {
		go takeSnapshot(&wg, subreddit["subreddit"].(string), geddit.PopularitySort(subreddit["sort"].(string)), ch,
			reddit)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	storeSnapshots(ch)
}

func takeSnapshot(wg *sync.WaitGroup, subreddit string, sort geddit.PopularitySort, ch chan catcher.SubredditSnapshot,
	reddit catcher.RedditClient) {
	defer wg.Done()
	snapshot := catcher.TakeSnapshot(reddit, subreddit, sort)
	ch <- snapshot
}

func storeSnapshots(ch chan catcher.SubredditSnapshot) {
	for msg := range ch {
		storer.StoreItem(msg, dbUrl, dbName, snapshotsCollection)
	}
}
