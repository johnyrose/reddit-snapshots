package storer

import (
	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func NewMongoClient(url string) mongoClient {
	client, err := mgo.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	return mongoClient{*client}
}

type mongoClient struct {
	Session mgo.Session
}

func (c mongoClient) insertToDatabase(dbName string, collectionName string, info bson.M) {
	client := c.Session
	collection := client.DB(dbName).C(collectionName)
	err := collection.Insert(info)
	if err != nil {
		log.Fatal(err)
	}
}
