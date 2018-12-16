package manager

import (
	"github.com/jzelinskie/reddit"
	"gopkg.in/mgo.v2/bson"
	"sync"

	"github.com/Ripolak/reddit-snapshots/catcher"
	"github.com/Ripolak/reddit-snapshots/storer"
)

func Entrypoint() {
	snapshotConfig := storer.LoadConfiguration(dbUrl, dbName, configCollection)
	subreddits := snapshotConfig.Subreddits
	fetchSnapshots(subreddits)
}

func fetchSnapshots(subreddits []bson.M) {
	var wg sync.WaitGroup
	wg.Add(len(subreddits))
	ch := make(chan catcher.SubredditSnapshot, len(subreddits))
	takeSnapshots(subreddits, &wg, ch)

	go func() {
		wg.Wait()
		close(ch)
	}()

	storeSnapshots(ch)
}

func takeSnapshots(subreddits []bson.M, wg *sync.WaitGroup,
	ch chan catcher.SubredditSnapshot) {
	for _, subreddit := range subreddits {
		go func(subreddit string, sort geddit.PopularitySort) {
			defer wg.Done()
			snapshot := catcher.TakeSnapshot(reddit, subreddit, sort)
			ch <- snapshot
		}(subreddit["subreddit"].(string), geddit.PopularitySort(subreddit["sort"].(string)))
	}
}

func storeSnapshots(ch chan catcher.SubredditSnapshot) {
	for msg := range ch {
		storer.StoreItem(msg, dbUrl, dbName, snapshotsCollection)
	}
}
