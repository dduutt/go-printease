package internal

import (
	"context"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
)

var prColl = mongox.NewCollection[PrintRecord](DBClient.NewDatabase("printease"), "print_records")

type PrintRecord struct {
	Model        `bson:"inline"`
	TemplateName string               `bson:"template_name" json:"template_name"`
	Fields       []*map[string]string `bson:"fields" json:"fields"`
	BatchCode    string               `bson:"batch_code" json:"batch_code"`
	Printer      string               `bson:"printer" json:"printer"`
}

func (p *PrintRecord) Create(pi *PrintRecord) error {
	pi.CreatedAt = pi.DefaultCreatedAt()
	_, err := prColl.Creator().InsertOne(context.Background(), pi)
	return err
}

func (p *PrintRecord) FindByBatchCode(batchCode string) ([]*PrintRecord, error) {
	f := query.Eq("batch_code", batchCode)

	return prColl.Finder().Filter(f).Find(context.Background())
}

func (p *PrintRecord) FindByBatchCodes(batchCodes []string) ([]*PrintRecord, error) {

	f := query.In("batch_code", batchCodes...)
	return prColl.Finder().Filter(f).Find(context.Background())

}

func (p *PrintRecord) CreateMany(pi []*PrintRecord) error {
	for _, v := range pi {
		v.CreatedAt = v.DefaultCreatedAt()
	}
	_, err := prColl.Creator().InsertMany(context.Background(), pi)
	return err
}
