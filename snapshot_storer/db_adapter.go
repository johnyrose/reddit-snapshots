package snapshot_storer

import (
	"github.com/globalsign/mgo"
	"log"
)

type MongoClient struct {
	Url string
}

func (c MongoClient) ConnectToDatabase() mgo.Session {
	client, err := mgo.Dial(c.Url)
	if err != nil {
		log.Fatal(err)
	}
	return *client
}
