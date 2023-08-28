package product

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/hope-ag/go-dynamo/core/entities"
	"github.com/hope-ag/go-dynamo/core/entities/product"
)

type Rules struct {}

func NewRules() *Rules {
	return &Rules{}
}

func (r *Rules) ConvertIoReaderToStruct(body io.Reader, model any) (any, error) {
	if body == nil {
		return nil, errors.New("body is required")
	}
	return model, json.NewDecoder(body).Decode(model)
}

func (r *Rules) GetMock() any {
	return product.Product {
		Base: entities.Base{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: uuid.New().String(),
	}
}

func (r *Rules) Validate(model any) error {
	productModel,err := product.InterfaceToModel(model)

	if err != nil {
		return err
	}
	return validation.ValidateStruct(productModel, 
		validation.Field(&productModel.ID, validation.Required, is.UUIDv4),
		validation.Field(&productModel.Name, validation.Required, validation.Length(3, 50)),
	)
}

func (r *Rules) Migrate(connection *dynamodb.DynamoDB) error {
	return r.createTable(connection)
}

func (r *Rules) createTable(connection *dynamodb.DynamoDB) error {
	table := &product.Product{}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("_id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("_id"),
				KeyType: aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits: aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(table.TableName()),
	}
	response,err := connection.CreateTable(input)
	if err != nil && strings.Contains(err.Error(), "Table already exists") {
		return nil
	}
	if response != nil && strings.Contains(response.GoString(), "CREATING") {
		time.Sleep(3 * time.Second)
		err = r.createTable(connection)
		if err != nil {
			return err
		}
	}
	return err
}
