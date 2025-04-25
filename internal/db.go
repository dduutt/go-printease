package internal

import (
	"context"
	"time"

	"github.com/chenmingyong0423/go-mongox/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const DBURL = "mongodb+srv://dote27:<pwd>@cluster0.uotcs1c.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
const PRDBURL = "mongodb://jldg:<pwd>@10.50.21.152:27017/printease?authSource=admin&tls=false"

var DBClient = InitDB()

type Model struct {
	ID        bson.ObjectID `bson:"_id,omitempty" mongox:"autoID" json:"id"`
	CreatedAt time.Time     `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updatedAt"`
}

// wails:ignore
func (m *Model) DefaultId() bson.ObjectID {
	if m.ID.IsZero() {
		m.ID = bson.NewObjectID()
	}
	return m.ID
}

// wails:ignore
func (m *Model) DefaultCreatedAt() time.Time {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().Local()
	}
	return m.CreatedAt
}

// wails:ignore
func (m *Model) DefaultUpdatedAt() time.Time {
	m.UpdatedAt = time.Now().Local()
	return m.UpdatedAt
}

func InitDB() *mongox.Client {
	options := options.Client().ApplyURI(PRDBURL)
	client, err := mongo.Connect(options)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		panic(err)
	}
	return mongox.NewClient(client, &mongox.Config{})
}

func Close() error {
	return DBClient.Disconnect(context.Background())
}
