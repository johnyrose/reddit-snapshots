package manager

import (
	"github.com/jzelinskie/reddit"
	"gopkg.in/mgo.v2/bson"
	"sync"

	"github.com/Ripolak/reddit-snapshots/catcher"
	"github.com/Ripolak/reddit-snapshots/storer"
)

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
