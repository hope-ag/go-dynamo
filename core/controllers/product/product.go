package product

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"github.com/hope-ag/go-dynamo/core/entities/product"
	"github.com/hope-ag/go-dynamo/core/repository/adapter"
)

type Interface interface {
	Create(body *product.Product) (uuid.UUID, error)
	ListAll()([]product.Product, error)
	ListOne(id uuid.UUID)(product.Product, error)
	Update(id uuid.UUID, product *product.Product) error
	Delete(id uuid.UUID) error
}

type Controller struct {
	repository adapter.Interface
}

func NewController(repository adapter.Interface) Interface {
	return &Controller{
		repository: repository,
	}
}


func (c *Controller) ListOne(id uuid.UUID) (product.Product, error) {
	entity := product.Product{}
	entity.ID = id
	response, err := c.repository.FindOne(entity.GetFilterId(), entity.TableName())
	if err != nil {
		return entity, err
	}
	return product.ParseDynamoAttributeToStruct(response.Item)
}

func (c *Controller) ListAll() ([]product.Product, error) {
	entities := []product.Product{}
	var entity product.Product

	filter := expression.Name("name").NotEqual(expression.Value(""))
	condition,err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return entities, err
	}
	response, err := c.repository.FindAll(condition, entity.TableName())
	if err != nil {
		return entities, err
	}
	if (response != nil) {
		for _,value := range response.Items {
			entity, err := product.ParseDynamoAttributeToStruct(value)
			if (err != nil) {
				return entities, err
			}
			entities = append(entities, entity)
		}
	} 
	return entities, nil
}

func (c *Controller) Create(entity *product.Product) (uuid.UUID, error) {
	entity.SetCreatedAt()
	entity.GenerateID()
	entity.SetUpdatedAt()
	_,err := c.repository.CreateOrUpdate(entity.GetMap(), entity.TableName())
	return entity.ID, err
}

func (c *Controller) Update(id uuid.UUID, product *product.Product) error {
	found,err := c.ListOne(id)
	if err != nil {
		return err
	}
	fmt.Println(product)
	found.ID = product.ID
	found.Name = product.Name
	found.UpdatedAt = time.Now()
	_,err = c.repository.CreateOrUpdate(found.GetMap(), product.TableName())
	return err
}

func (c *Controller) Delete(id uuid.UUID) error {
	entity,err := c.ListOne(id)
	if err != nil {
		return err
	}
	_,err = c.repository.Delete(entity.GetFilterId(), entity.TableName())
	return err
}