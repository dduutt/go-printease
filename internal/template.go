package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/bsonx"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var tpColl = mongox.NewCollection[Template](DBClient.NewDatabase("printease"), "templates")

type Template struct {
	Model       `bson:"inline"`
	Name        string              `bson:"name" json:"name"`
	Path        string              `bson:"path" json:"path"`
	Description string              `bson:"description" json:"description"`
	Disabled    int                 `bson:"disabled" json:"disabled"`
	Filds       []map[string]string `bson:"fields" json:"fields"`
	Datas       []map[string]string `bson:"datas" json:"datas"`
}

type ListByNameResp struct {
	Total int64       `json:"total"`
	List  []*Template `json:"list"`
}

func (t *Template) Create(ti *Template) error {
	ti.CreatedAt = t.DefaultCreatedAt()
	err := ti.readFromXlsx()
	if err != nil {
		return fmt.Errorf("read from xlsx error:%s", err)
	}
	_, err = tpColl.Creator().InsertOne(context.Background(), ti)
	return err
}

// 根据名称模糊查询
func (t *Template) ListByName(name string, skip, limit int) (*ListByNameResp, error) {
	total, err := t.Count(name)
	if err != nil {
		return nil, err
	}
	tps, err := tpColl.Finder().Filter(query.Regex("name", name)).Sort(bsonx.M("created_at", -1)).Skip(int64(skip)).Limit(int64(limit)).Find(context.Background(), options.Find().SetProjection(bsonx.M("datas", 0)))
	if err != nil {
		return nil, err
	}
	return &ListByNameResp{Total: total, List: tps}, nil
}

func (t *Template) ListUsedByName(name string) (*ListByNameResp, error) {
	total, err := t.Count(name)
	if err != nil {
		return nil, err
	}
	tps, err := tpColl.Finder().Filter(query.Regex("name", name)).Sort(bsonx.M("created_at", -1)).Find(context.Background(), options.Find().SetProjection(bsonx.M("datas", 0)))
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
	if r.DeletedCount == 0 {
		return fmt.Errorf("delete template error: %v", err)
	}
	return nil
}

func (t *Template) FindByName(name string) (*Template, error) {
	return tpColl.Finder().Filter(query.Eq("name", name)).FindOne(context.Background(), options.FindOne().SetProjection(bsonx.M("datas", 0)))
}

func (t *Template) Count(name string) (int64, error) {
	return tpColl.Finder().Filter(query.Regex("name", name)).Count(context.Background())
}

func (t *Template) readFromXlsx() error {
	f, err := excelize.OpenFile(t.Path)
	if err != nil {
		return err
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Println("close file error:", err)
		}
	}()

	sheetName := f.GetSheetList()[0]
	if sheetName == "" {
		return fmt.Errorf("sheet name is empty")
	}
	rows, err := f.Rows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to read rows: %v", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("close rows error:", err)
		}
	}()
	var headers []string
	if rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("failed to read header row: %v", err)
		}
		headers = row
		for _, value := range row {
			m := make(map[string]string)
			m["name"] = value
			m["key"] = value
			m["value"] = ""
			t.Filds = append(t.Filds, m)
		}
	}
	rowIdx := 1
	for rows.Next() {
		rowIdx++
		row, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("failed to read row %d: %v", rowIdx, err)
		}
		m := make(map[string]string)

		for i, value := range row {
			if i < len(headers) {
				m[headers[i]] = value
			}
		}
		t.Datas = append(t.Datas, m)
	}
	return nil
}

func (t *Template) FindDatasByKeys(id string, keys []map[string]string) ([]map[string]string, error) {
	// 转换字符串ID为MongoDB ObjectID
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID格式: %v", err)
	}

	// 构建动态过滤条件
	conditions := make([]bson.M, 0, len(keys))
	for _, key := range keys {
		if key["key"] == "" || key["value"] == "" {
			continue
		}
		conditions = append(conditions, bson.M{
			"$eq": []any{
				fmt.Sprintf("$$this.%s", key["key"]),
				key["value"],
			},
		})
	}

	// 如果没有有效条件则返回空
	// if len(conditions) == 0 {
	// 	return []map[string]string{}, nil
	// }

	// 构建聚合管道
	pipeline := []bson.M{
		{"$match": bson.M{"_id": objID}},
		{"$project": bson.M{
			"filteredDatas": bson.M{
				"$filter": bson.M{
					"input": "$datas",
					"as":    "this",
					"cond":  bson.M{"$and": conditions},
				},
			},
			"_id": 0,
		}},
	}

	// 执行聚合查询
	cursor, err := tpColl.Collection().Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("聚合查询失败: %v", err)
	}
	defer func() {
		if err := cursor.Close(context.Background()); err != nil {
			log.Printf("关闭游标失败: %v", err)
		}
	}()

	// 解析查询结果
	var result struct {
		FilteredDatas []map[string]string `bson:"filteredDatas"`
	}
	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("结果解析失败: %v", err)
		}
	}

	return result.FilteredDatas, nil
}
