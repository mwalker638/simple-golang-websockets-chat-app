package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// structure representing the connections table
type connection struct {
	ConnectionID string
}

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch route := request.RequestContext.RouteKey; route {
	case "$connect":
		return doConnect(ctx, request)
	case "sendmessage":
		return doSendmessage(ctx, request)
	case "$disconnect":
		return doDisconnect(ctx, request)
	default:
		return handleError("unexepcted route: " + route)
	}
}

func doConnect(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	label := "doConnect: "
	cID := request.RequestContext.ConnectionID
	log.Println(label, "cID: ", cID)

	session, err := session.NewSession()
	if err != nil {
		return handleError("failed to establish aws session: " + err.Error())
	}
	dbSvc := dynamodb.New(session)

	av, err := dynamodbattribute.MarshalMap(&connection{ConnectionID: cID})
	if err != nil {
		return handleError(label + "failed to marshal item: " + err.Error())
	}

	if _, err = dbSvc.PutItem(&dynamodb.PutItemInput{Item: av, TableName: aws.String(os.Getenv("TABLE_NAME"))}); err != nil {
		return handleError(label + "connectionid: putitem failed: " + err.Error())
	}
	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func doSendmessage(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	label := "doSendmessage: "
	session, err := session.NewSession()
	if err != nil {
		return handleError(label + "failed to establish aws session: " + err.Error())
	}

	log.Println(label, "body: ", request.Body)

	type message struct {
		Message string
		Data    string
	}
	msg := message{}
	if err = json.Unmarshal([]byte(request.Body), &msg); err != nil {
		return handleError(label + "failed to unmarshal body: " + err.Error())
	}

	endurl := "https://" + request.RequestContext.DomainName + "/" + request.RequestContext.Stage
	apigw := apigatewaymanagementapi.New(session, aws.NewConfig().WithEndpoint(endurl))

	dbSvc := dynamodb.New(session)
	result, err := dbSvc.Scan(&dynamodb.ScanInput{TableName: aws.String(os.Getenv("TABLE_NAME"))})
	if err != nil {
		return handleError(label + "dynamodb.Scan(ConnectionID) failed: " + err.Error())
	}

	conList := make([]connection, *result.Count)
	if err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &conList); err != nil {
		return handleError(label + "UnmarshalListOfMaps() failed: " + err.Error())
	}

	for _, c := range conList {
		log.Println(label, "posting to cid: ", c.ConnectionID)
		_, err = apigw.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{ConnectionId: &c.ConnectionID, Data: []byte(msg.Data)})
		if err != nil {
			aerr, ok := err.(awserr.Error)
			if ok && aerr.Code() == apigatewaymanagementapi.ErrCodeGoneException {
				if err = deleteConnection(c.ConnectionID); err != nil {
					return handleError(label + "deleteConnection() failed: " + err.Error())
				}
			} else {
				log.Println(label, "Error: PostToConnection() failed: ", err.Error())
			}
		}
	}
	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func doDisconnect(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	if err := deleteConnection(request.RequestContext.ConnectionID); err != nil {
		return handleError("doDisconnect: deleteConnection() failed: " + err.Error())
	}
	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

//
// Local helper functions
//

// General error logger and return payload builder
func handleError(message string) (events.APIGatewayProxyResponse, error) {
	emsg := "Error: " + message
	log.Println(emsg)
	return events.APIGatewayProxyResponse{Body: emsg, StatusCode: 500}, nil
}

// Remove connections from the connections_table.
func deleteConnection(cID string) error {
	log.Println("deleteConnection: cID=", cID)
	session, err := session.NewSession()
	if err != nil {
		return err
	}

	dbSvc := dynamodb.New(session)
	delIn := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionID": {S: aws.String(cID)},
		},
		ReturnValues: aws.String("ALL_OLD"),
		TableName:    aws.String(os.Getenv("TABLE_NAME")),
	}
	if _, err := dbSvc.DeleteItem(delIn); err != nil {
		return err
	}
	return nil
}

//
// Main
//
func main() {
	lambda.Start(handleRequest)
}
