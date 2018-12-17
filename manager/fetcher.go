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

type config struct {
	DbConfig     dbConfig
	RedditConfig redditConfig
}

func (c config) ProcessConfig() config {
	err := envconfig.Process("", &c.DbConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = envconfig.Process("", &c.RedditConfig)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func (c config) GenerateReddit() catcher.RedditClient {
	reddit := catcher.RedditClient{
		ClientID:     c.RedditConfig.ClientID,
		ClientSecret: c.RedditConfig.ClientSecret,
		Username:     c.RedditConfig.Username,
		Password:     c.RedditConfig.Password,
	}
	return reddit
}

func (c config) GenerateStorer() storer.SnapshotStorer {
	mongoClient := storer.NewMongoClient(c.DbConfig.DbUrl)
	snapshotsStorer := storer.DatabaseStorer{
		MongoClient: mongoClient,
		DbName:      c.DbConfig.DbName,
		Collection:  c.DbConfig.SnapshotsCollection,
	}
	return snapshotsStorer
}

func Entrypoint() {
	var c config
	c = c.ProcessConfig()
	subredditStorer := c.GenerateStorer()
	reddit := c.GenerateReddit()
	snapshotConfig := storer.LoadConfiguration(c.DbConfig.DbUrl, c.DbConfig.DbName, c.DbConfig.ConfigCollection)
	subreddits := snapshotConfig.Subreddits
	fetchSnapshots(subreddits, reddit, subredditStorer)
}

func fetchSnapshots(subreddits []bson.M, reddit catcher.RedditClient, subredditStorer storer.SnapshotStorer) {
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

	wg.Add(len(subreddits))
	storeSnapshots(&wg, ch, subredditStorer)
	wg.Wait()
}

func takeSnapshot(wg *sync.WaitGroup, subreddit string, sort geddit.PopularitySort, ch chan catcher.SubredditSnapshot,
	reddit catcher.RedditClient) {
	defer wg.Done()
	snapshot := catcher.TakeSnapshot(reddit, subreddit, sort)
	ch <- snapshot
}

func storeSnapshots(wg *sync.WaitGroup, ch chan catcher.SubredditSnapshot, snapshotStorer storer.SnapshotStorer) {
	for msg := range ch {
		go func() {
			defer wg.Done()
			snapshotStorer.StoreItem(msg)
		}()
	}
}
