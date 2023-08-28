package product

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/hope-ag/go-dynamo/core/entities"
)

type Product struct {
	entities.Base
	Name string `json:"name"`
}

func InterfaceToModel(data any) (*Product, error) {
	instance := Product{}
	bytes, err := json.Marshal(data)

	if err != nil {
		return &instance, err
	}
	return &instance, json.Unmarshal(bytes, &instance)
}

func (p *Product) TableName() string {
	return "products"
}

func (p *Product) GetMap() map[string]any {
	return map[string]any{
		"_id":       p.ID.String(),
		"name":      p.Name,
		"createdAt": p.CreatedAt.Format(entities.GetTimeFormat()),
		"updatedAt": p.UpdatedAt.Format(entities.GetTimeFormat()),
	}
}

func (p *Product) GetFilterId() map[string]any {
	return map[string]any{"_id": p.ID.String()}
}

func (p *Product) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func ParseDynamoAttributeToStruct(attribute map[string]*dynamodb.AttributeValue) (Product, error) {
	p := Product{}
	var err error
	if attribute == nil || (attribute != nil && len(attribute) == 0) {
		return p, errors.New("item not found")
	}

	for key, value := range attribute {
		if key == "_id" {
			p.ID, err = uuid.Parse(*value.S)
			if p.ID == uuid.Nil {
				err = errors.New("item not found")
			}
		}
		if key == "name" {
			p.Name = *value.S
		}
		if key == "createdAt" {
			p.CreatedAt, err = time.Parse(entities.GetTimeFormat(), *value.S)
		}
		if key == "updatedAt" {
			p.UpdatedAt, err = time.Parse(entities.GetTimeFormat(), *value.S)
		}
		if err != nil {
			return p, err
		}
	}
	return p, err
}
