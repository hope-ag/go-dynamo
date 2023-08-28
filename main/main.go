package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hope-ag/go-dynamo/config"
	"github.com/hope-ag/go-dynamo/core/repository/adapter"
	"github.com/hope-ag/go-dynamo/core/repository/instance"
	"github.com/hope-ag/go-dynamo/core/routes"
	"github.com/hope-ag/go-dynamo/core/rules"
	RulesProduct "github.com/hope-ag/go-dynamo/core/rules/product"
	"github.com/hope-ag/go-dynamo/utils/env"
	"github.com/hope-ag/go-dynamo/utils/logger"
)

func init() {
	env.InitEnvironment()
}

func main() {
	config := config.GetConfig()
	connection := instance.GetConnection()
	repository := adapter.NewAdapter(connection)

	logger.INFO("----------\n\nStarting Service\n\n----------", nil)

	err := Migrate(connection)
	if len(err) > 0 {
		for _, e := range err {
			logger.PANIC("Migration error", e)
		}
	}

	logger.PANIC("", CheckTables(connection))

	port := fmt.Sprintf(":%v", config.Port)
	router := routes.NewRouter().SetRouters(repository)
	logger.INFO("App is running on port ", port)

	server := http.ListenAndServe(port, router)
	log.Fatal(server)
}

func Migrate(connection *dynamodb.DynamoDB) []error {
	var errors []error
	callMigrateAndAppendErrors(&errors, connection, &RulesProduct.Rules{})
	return errors
}

func callMigrateAndAppendErrors(errors *[]error, connection *dynamodb.DynamoDB, rule rules.Interface) {
	err := rule.Migrate(connection)
	if err != nil {
		*errors = append(*errors, err)
	}
}

func CheckTables(connection *dynamodb.DynamoDB) error {
	response, err := connection.ListTables(&dynamodb.ListTablesInput{})

	if err != nil {
		return err
	}
	if len(response.TableNames) == 0 {
		logger.INFO("No tables found", nil)
	} else {
		for _, name := range response.TableNames {
			logger.INFO("Table Name: ", *name)
		}
	}
	return nil
}
