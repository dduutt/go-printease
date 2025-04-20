package internal

import (
	"context"
	"fmt"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/bsonx"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var tpColl = mongox.NewCollection[Template](DBClient.NewDatabase("printease"), "templates")

type Template struct {
	Model       `bson:"inline"`
	Name        string `bson:"name" json:"name"`
	Path        string `bson:"path" json:"path"`
	Description string `bson:"description" json:"description"`
	InUse       int    `bson:"in_use" json:"inUse"`
}

type ListByNameResp struct {
	Total int64       `json:"total"`
	List  []*Template `json:"list"`
}

func (t *Template) Create(ti *Template) error {
	ti.CreatedAt = t.DefaultCreatedAt()
	r, err := tpColl.Creator().InsertOne(context.Background(), ti)
	if err != nil {
		return err
	}
	objID, ok := r.InsertedID.(bson.ObjectID)
	if !ok {
		return fmt.Errorf("insert template error: %v", err)
	}
	if objID.IsZero() {
		return fmt.Errorf("insert template error: %v", err)
	}
	_, errs := LoadTemplateDatas(objID.Hex(), ti.Path)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println("insert template data error:", err)
		}
		return fmt.Errorf("insert template data error: %v", errs)
	}
	return nil
}

// 根据名称模糊查询
func (t *Template) ListByName(name string, skip, limit int) (*ListByNameResp, error) {
	total, err := t.Count(name)
	if err != nil {
		return nil, err
	}
	tps, err := tpColl.Finder().Filter(query.Regex("name", name)).Sort(bsonx.M("created_at", -1)).Skip(int64(skip)).Limit(int64(limit)).Find(context.Background())
	if err != nil {
		return nil, err
	}
	return &ListByNameResp{Total: total, List: tps}, nil
}

func (t *Template) Update(ut Template) error {
	ut.UpdatedAt = t.DefaultUpdatedAt()
	_, err := tpColl.Updater().Filter(query.Id(ut.ID)).
		Updates(update.SetFields(ut)).UpdateOne(context.Background())
	return err
}

func (t *Template) Delete(id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	r, err := tpColl.Deleter().Filter(query.Id(objID)).DeleteOne(context.Background())
	fmt.Println("delete result:", r.DeletedCount)
	return err
}

func (t *Template) FindByName(name string) (*Template, error) {
	return tpColl.Finder().Filter(query.Eq("name", name)).FindOne(context.Background())
}

func (t *Template) Count(name string) (int64, error) {
	return tpColl.Finder().Filter(query.Regex("name", name)).Count(context.Background())
}
