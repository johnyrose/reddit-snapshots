package reddit_snapshot_catcher

import (
	"log"

	"github.com/jzelinskie/reddit"
)

type RedditClient struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

func (c RedditClient) generateSession() *geddit.OAuthSession {
	o, err := geddit.NewOAuthSession(
		c.ClientID,
		c.ClientSecret,
		"test",
		"http://redirect.url",
	)
	if err != nil {
		log.Fatal(err)
	}
	err = o.LoginAuth(c.Username, c.Password)
	if err != nil {
		log.Fatal(err)
	}
	return o
}
