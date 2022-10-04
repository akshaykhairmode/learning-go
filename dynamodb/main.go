package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//AWS Session object
var awsSession = session.Must(session.NewSession(
	&aws.Config{
		Region: aws.String("ap-south-1"),
	},
))

//DynamoDB object
var ddb = dynamodb.New(awsSession)

func main() {

	log.SetOutput(os.Stdout)

	listTables()

	employeeTableName := "employee"

	// deleteTable(employeeTableName)

	// createTable(employeeTableName)

	//inserts
	createItem(employeeTableName, "1", "John Doe", "25")
	createItem(employeeTableName, "2", "Jane Doe", "20")

	//selects
	getItem(employeeTableName, "1")
	getItem(employeeTableName, "2")

	//deletes
	deleteItem(employeeTableName, "1")
	getItem(employeeTableName, "1")

	//updates
	createItem(employeeTableName, "2", "Jane Doe Junior", "20")
	getItem(employeeTableName, "2")

}

func listTables() {

	input := &dynamodb.ListTablesInput{
		ExclusiveStartTableName: nil,
		Limit:                   aws.Int64(10),
	}

	output, err := ddb.ListTables(input)
	if err != nil {
		log.Println("error while listing tables", err)
		return
	}

	log.Println("Printing Dynamo DB Tables")

	for _, tableName := range output.TableNames {
		log.Println("*", *tableName)
	}

}

func createTable(tableName string) {

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{ //We mention the key and its data type here.
				AttributeName: aws.String("empId"),
				AttributeType: aws.String("N"),
			},
		},
		//Mention if our attribute is partition key or sort key.
		//We use hash for partition key.
		KeySchema: []*dynamodb.KeySchemaElement{{
			AttributeName: aws.String("empId"),
			KeyType:       aws.String("HASH"),
		}},
		TableName:   aws.String(tableName),
		BillingMode: aws.String("PAY_PER_REQUEST"),
	}

	//create table is async, aws will give response early but may take some time to create the table.
	output, err := ddb.CreateTable(input)
	if err != nil {
		log.Println("error while creating table", err)
		return
	}

	log.Printf("Create table successful with name : %s , output is : %s", tableName, output)
}

func deleteTable(tableName string) {

	input := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	}

	//delete table is async, aws will give response early but may take some time to delete the table.
	output, err := ddb.DeleteTable(input)
	if err != nil {
		log.Println("error while deleting table", err)
		return
	}

	log.Printf("Delete table successful with name : %s, output is : %s", tableName, output)

}

func createItem(tableName, empId, name, age string) {

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"empId": {
				N: &empId, //N is number.
			},
			"name": {
				S: &name, //S is string
			},
			"age": {
				N: &age,
			},
			"active": {
				BOOL: aws.Bool(true),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := ddb.PutItem(input)
	if err != nil {
		log.Println("error while putting item", err)
		return
	}

	log.Printf("Put item success for empId : %s", empId)

}

func getItem(tableName, empId string) {

	input := &dynamodb.GetItemInput{
		ConsistentRead: aws.Bool(true),
		Key: map[string]*dynamodb.AttributeValue{
			"empId": {
				N: &empId,
			},
		},
		TableName: aws.String(tableName),
	}

	output, err := ddb.GetItem(input)
	if err != nil {
		log.Println("error while getting item", err)
		return
	}

	//We get the data in the Item field. Item field is a map. In case the item does not exist the map is nil.
	if output.Item == nil {
		log.Printf("Item does not exist for empId : %s", empId)
		return
	}

	name := *output.Item["name"].S
	age := *output.Item["age"].N
	active := *output.Item["active"].BOOL

	log.Printf("Get item success, data :- empId : %s, name : %s, age : %s, active : %v", empId, name, age, active)
}

func deleteItem(tableName, empId string) {

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"empId": {
				N: &empId,
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := ddb.DeleteItem(input)
	if err != nil {
		log.Println("error while deleting item", err)
		return
	}

	log.Printf("Delete item success for empId : %s", empId)
}
