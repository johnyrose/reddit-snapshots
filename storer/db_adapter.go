package storer

import (
	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
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

func (c MongoClient) insertToDatabase(dbName string, collectionName string, info bson.M) {
	client := c.ConnectToDatabase()
	collection := client.DB(dbName).C(collectionName)
	err := collection.Insert(info)
	if err != nil {
		log.Fatal(err)
	}
}
