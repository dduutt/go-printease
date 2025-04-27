package internal

import (
	"context"
	_ "embed"
	"log"
	"time"

	"github.com/chenmingyong0423/go-mongox/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

//go:embed ".env"
var DB_URI string
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
		m.CreatedAt = time.Now()
	}
	return m.CreatedAt
}

// wails:ignore
func (m *Model) DefaultUpdatedAt() time.Time {
	m.UpdatedAt = time.Now()
	return m.UpdatedAt
}

func InitDB() *mongox.Client {

	options := options.Client().ApplyURI(DB_URI)
	client, err := mongo.Connect(options)
	if err != nil {
		log.Printf("连接MongoDB失败: %v\n", err)
		panic(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Printf("MongoDB ping失败: %v\n", err)
		panic(err)
	}
	log.Println("成功连接到MongoDB")
	return mongox.NewClient(client, &mongox.Config{})
}

func Close() error {
	return DBClient.Disconnect(context.Background())
}
