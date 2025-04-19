package internal

import (
	"context"

	"github.com/chenmingyong0423/go-mongox/v2"
)

var tpdColl = mongox.NewCollection[TemplateData](DBClient.NewDatabase("printease"), "template_datas")

type TemplateData struct {
	Model        `bson:"inline"`
	TemplateName string           `bson:"template_name" json:"templateName"`
	TemplateID   string           `bson:"template_id" json:"templateId"`
	Data         []map[string]any `bson:"data" json:"data"`
}

func (t *TemplateData) Create(tdi *TemplateData) error {

	tdi.CreatedAt = t.defaultCreatedAt()
	_, err := tpdColl.Creator().InsertOne(context.Background(), tdi)
	return err
}
