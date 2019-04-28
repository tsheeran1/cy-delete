package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Age    int `json:"age"`
	Height int `json:"height"`
	Income int `json:"income"`
}

func handler(ctx context.Context, e Event) (string, error) {

	fmt.Println("Event:", e)
	fmt.Println("Age: ", e.Age)
	eage := e.Age
	return fmt.Sprintf("%v", eage*2), nil
}

func main() {
	lambda.Start(handler)
}
