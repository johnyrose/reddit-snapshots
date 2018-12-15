package snapshots_manager

import (
	"gopkg.in/mgo.v2/bson"
	"sync"

	"github.com/Ripolak/reddit-snapshots/reddit_snapshot_catcher"
	"github.com/Ripolak/reddit-snapshots/snapshot_storer"
)

func Entrypoint() {
	snapshotConfig := snapshot_storer.LoadConfiguration(dbUrl, dbName, configCollection)
	subreddits := snapshotConfig.Subreddits
	fetchSnapshots(subreddits)
}

func fetchSnapshots(subreddits []bson.M) {
	var wg sync.WaitGroup
	wg.Add(len(subreddits))
	ch := make(chan reddit_snapshot_catcher.SubredditSnapshot, len(subreddits))
	takeSnapshots(subreddits, &wg, ch)
	wg.Wait()
	close(ch)
	wg.Add(len(subreddits))
	storeSnapshots(ch, &wg)
	wg.Wait()
}

func takeSnapshots(subreddits []bson.M, wg *sync.WaitGroup,
	ch chan reddit_snapshot_catcher.SubredditSnapshot) {
	for _, subreddit := range subreddits {
		go func(subreddit string) {
			defer wg.Done()
			snapshot := reddit_snapshot_catcher.TakeSnapshot(reddit, subreddit, "hot")
			ch <- snapshot
		}(subreddit["subreddit"].(string))
	}
}

func storeSnapshots(ch chan reddit_snapshot_catcher.SubredditSnapshot, wg *sync.WaitGroup) {
	for msg := range ch {
		go func(snap reddit_snapshot_catcher.SubredditSnapshot) {
			defer wg.Done()
			snapshot_storer.StoreItem(snap, dbUrl, dbName, snapshotsCollection)
		}(msg)
	}
}
