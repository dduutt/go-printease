package internal

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Template struct {
	ID          string `bson:"id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Path        string `bson:"path" json:"path"`
	Description string `bson:"description" json:"description"`
	InUse       int    `bson:"inUse" json:"inUse"`
	CreatedAt   string `bson:"createdAt" json:"createdAt"`
	UpdatedAt   string `bson:"updatedAt" json:"updatedAt"`
}
type TemplateData struct {
	TemplateId string   `bson:"templateId" json:"templateId"`
	Fields     []string `bson:"fields" json:"fields"`
	Data       []any    `bson:"data" json:"data"`
}

func (t *Template) InsertOne(it Template) (Template, error) {
	c := MongoDB.Database("printease").Collection("templates")
	f := bson.D{{Key: "name", Value: it.Name}}
	rt := Template{}
	err := c.FindOne(context.TODO(), f).Decode(rt)
	if err == nil {
		return rt, fmt.Errorf("%s already exists", it.Name)
	}
	if err == mongo.ErrNoDocuments {
		it.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
		it.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		ir, e := c.InsertOne(context.TODO(), it)
		if e != nil {
			return rt, fmt.Errorf("插入失败: %v", e)
		}
		it.ID = ir.InsertedID.(string)
		return it, nil
	}
	return it, fmt.Errorf("查询失败: %v", err)
}

func (t *Template) QueryByName(name string) ([]*Template, error) {

	c := MongoDB.Database("printease").Collection("templates")
	f := bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: name}, {Key: "$options", Value: "i"}}}}
	cur, err := c.Find(context.TODO(), f)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	defer cur.Close(context.TODO())
	var results []*Template
	err = cur.All(context.TODO(), results)
	return results, err
}

func (t *Template) List(q string, offset, limit int64) ([]Template, error) {

	c := MongoDB.Database("printease").Collection("templates")
	f := bson.D{
		{Key: "name", Value: bson.D{{Key: "$regex", Value: q}, {Key: "$options", Value: "i"}}},
	}
	opts := options.Find().SetLimit(limit).SetSkip(offset)
	cur, err := c.Find(context.TODO(), f, opts)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	defer cur.Close(context.TODO())
	var results []Template
	err = cur.All(context.TODO(), &results)
	return results, err
}
