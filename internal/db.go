package internal

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const DBURL = "mongodb://jldg:123456@10.50.21.152:27017/printease?authSource=admin&tls=false"

var MongoDB *mongo.Client

func InitDB() *mongo.Client {
	options := options.Client().ApplyURI(DBURL)
	client, err := mongo.Connect(options)
	if err != nil {
		panic(err)
	}
	MongoDB = client
	return client
}
