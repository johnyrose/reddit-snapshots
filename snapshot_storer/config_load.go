package snapshot_storer

import (
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Config struct {
	ID         bson.ObjectId
	Subreddits []bson.M
}

func loadConfigurationFromDB(c MongoClient, dbName string, collectionName string) Config {
	client := c.ConnectToDatabase()
	collection := client.DB(dbName).C(collectionName)
	var result Config
	err := collection.Find(nil).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func LoadConfiguration(dbUrl string, dbName string, configCollection string) Config {
	client := MongoClient{
		Url: dbUrl,
	}
	result := loadConfigurationFromDB(client, dbName, configCollection)
	return result
}
