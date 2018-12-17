package catcher

import (
	"log"

	"github.com/jzelinskie/reddit"
)

type RedditClient struct {
	Session *geddit.OAuthSession
}

func NewRedditClient(clientID string, clientSecret string, username string, password string) RedditClient {
	o, err := geddit.NewOAuthSession(
		clientID,
		clientSecret,
		"test",
		"http://redirect.url",
	)
	if err != nil {
		log.Fatal(err)
	}
	err = o.LoginAuth(username, password)
	if err != nil {
		log.Fatal(err)
	}
	return RedditClient{o}
}
