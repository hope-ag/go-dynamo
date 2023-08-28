package rules

import (
	"io"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Interface interface {
	Migrate(connection *dynamodb.DynamoDB) error
	ConvertIoReaderToStruct(body io.Reader, model any) (any, error)
	GetMock() any
	Validate(model any) error
}
