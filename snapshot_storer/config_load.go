package snapshot_storer

import "gopkg.in/mgo.v2/bson"

type Config struct {
	ID         bson.ObjectId
	Subreddits []bson.M
}
