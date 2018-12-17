package storer

import (
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Config struct {
	ID         bson.ObjectId
	Subreddits []bson.M
}

func loadConfigurationFromDB(c mongoClient, dbName string, collectionName string) Config {
	client := c.Session
	collection := client.DB(dbName).C(collectionName)
	var result Config
	err := collection.Find(nil).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func LoadConfiguration(dbUrl string, dbName string, configCollection string) Config {
	client := NewMongoClient(dbUrl)
	result := loadConfigurationFromDB(client, dbName, configCollection)
	return result
}
