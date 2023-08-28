package adapter

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type Database struct {
	connection *dynamodb.DynamoDB
	logMode bool
}

type Interface interface {
	Health() bool
	FindOne(condition map[string]any, tableName string) (*dynamodb.GetItemOutput, error)
	FindAll(condition expression.Expression, tableName string) (*dynamodb.ScanOutput, error)
	CreateOrUpdate(entity any, tableName string) (*dynamodb.PutItemOutput, error)
	Delete(entity map[string] any, tableName string) (*dynamodb.DeleteItemOutput, error)
}

func NewAdapter(con *dynamodb.DynamoDB) Interface {
	return & Database {
		connection: con,
		logMode: false,
	}
}

func (db *Database) Health() bool {
	_,err := db.connection.ListTables(&dynamodb.ListTablesInput{})
	return err == nil
}

func (db *Database) FindAll(condition expression.Expression, tableName string) (*dynamodb.ScanOutput, error) {

	input := &dynamodb.ScanInput {
		ExpressionAttributeNames: condition.Names(),
		ExpressionAttributeValues: condition.Values(),
		FilterExpression: condition.Filter(),
		ProjectionExpression: condition.Projection(),
		TableName: aws.String(tableName),
	}
	return db.connection.Scan(input)
}

func (db *Database) FindOne(condition map[string]any, tableName string) (*dynamodb.GetItemOutput, error) {
	conditionParsed,err := dynamodbattribute.MarshalMap(condition)
	if (err != nil) {
		return nil, err
	}
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: conditionParsed,
	}
	return db.connection.GetItem(input)
}

func (db *Database) CreateOrUpdate(entity any, tableName string) (*dynamodb.PutItemOutput, error) {
	entityParsed, err := dynamodbattribute.MarshalMap(entity)
	if (err != nil) {
		return nil, err
	}
  input := &dynamodb.PutItemInput{
		Item: entityParsed,
		TableName: aws.String(tableName),
	}
	return db.connection.PutItem(input)
}


func (db *Database) Delete(entity map[string] any, tableName string) (*dynamodb.DeleteItemOutput, error) {
	entityParsed, err := dynamodbattribute.MarshalMap(entity)
	if (err != nil) {
		return nil, err
	}
	input := &dynamodb.DeleteItemInput{
		Key: entityParsed,
		TableName: aws.String(tableName),
	}
	return db.connection.DeleteItem(input)
}