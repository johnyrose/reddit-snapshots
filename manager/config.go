package manager

import (
	"github.com/Ripolak/reddit-snapshots/catcher"
	"github.com/Ripolak/reddit-snapshots/storer"
	"github.com/kelseyhightower/envconfig"
	"log"
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
	reddit := catcher.NewRedditClient(c.RedditConfig.ClientID, c.RedditConfig.ClientSecret,
		c.RedditConfig.Username, c.RedditConfig.Password)
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
