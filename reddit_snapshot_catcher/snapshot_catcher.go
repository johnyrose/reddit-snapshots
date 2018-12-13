package reddit_snapshot_catcher

import "github.com/jzelinskie/reddit"

type SubredditSnapshot struct {
	Subreddit string
	Time      string
	Posts     []*geddit.Submission
}
