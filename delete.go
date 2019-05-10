package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	AccessToken string `json:"accessToken"`
}

type Record struct {
	Userid string `json:"Userid"`
	Age    int    `json:"age"`
	Height int    `json:"height"`
	Income int    `json:"income"`
}

type Keystruct struct {
	Userid string
}

func handler(ctx context.Context, e Event) error {

	fmt.Println("Event: ", e)
	// define AWS config, session and dynamodb objects
	config := &aws.Config{
		Region: aws.String("us-east-2"),
	}
	sess := session.Must(session.NewSession(config))
	dbc := dynamodb.New(sess)
	// create cognito service for this session
	cisp := cognitoidentityprovider.New(sess)

	accessToken := e.AccessToken

	// get userID for this user
	// define get user input struct
	getui := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	}
	// get the User object in to getuo
	getuo, err := cisp.GetUser(getui)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// userID is in getuo.UserAttributes[0].Value
	userID := getuo.UserAttributes[0].Value

	// Now define a key structure and dynamodbattribute Marshall it
	keyval := Keystruct{Userid: *userID}
	av, err := dynamodbattribute.MarshalMap(keyval)
	if err != nil {
		fmt.Println("Unable to marshal key structure")
		return err
	}
	// Now create a DeleteItemInput structure
	di := &dynamodb.DeleteItemInput{
		TableName: aws.String("compare-yourself"),
		Key:       av,
	}

	// Now delete the item with that key
	_, err = dbc.DeleteItem(di)
	if err != nil {
		fmt.Println("DeleteItem failure")
		return err
	}

	// ONLY IF WE SET THE RIGHT PARAMS in the av structure which we did not then dout.Attributes would contain a map[string]*AttributeValue which we can dynamodbattribute unmarshal into the record that was deleted

	// var r Record // record to unmarshall into
	// err = dynamodbattribute.UnmarshalMap(dout.Attributes, &r)
	// if err != nil {
	// 	fmt.Println("unable to unmarshal deleted record")
	// 	return r, err
	// }
	// fmt.Println(r)

	return nil
}

func main() {
	lambda.Start(handler)
}
