package main

/*
This lambda function will consume data from kinesis steam, extract the data and upload it to s3.
It will also send the data and event ID to DynamoDB.
The message attribute corresponding to "env" will be displayed in stdout.
The dynamodb table name and s3 bucket names are configured as environment variables in Lambda.
*/

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	lambda.Start(kinesisHandler)
}

func kinesisHandler(ctx context.Context, kinesisEvent events.KinesisEvent) error {
	data := kinesisEvent.Records[0]
	fmt.Println(data.EventID)

	c1 := make(chan *s3.PutObjectOutput)
	go uploadS3Bucket(string(data.Kinesis.Data), c1)

	for c1Msg := range c1 {
		fmt.Println(c1Msg)
	}

	c2 := make(chan *dynamodb.PutItemOutput)
	go sendToDynamo(data.EventID, string(data.Kinesis.Data), c2)

	for c2Msg := range c2 {
		fmt.Println(c2Msg)
	}

	return nil
}

func uploadS3Bucket(s string, c chan *s3.PutObjectOutput) {
	BUCKETNAME := os.Getenv("BUCKET")

	// Date for S3 filename
	now := time.Now()
	year := strconv.Itoa(int(now.Year()))
	month := strconv.Itoa(int(now.Month()))
	day := strconv.Itoa(int(now.Day()))
	hour := strconv.Itoa(int(now.Hour()))
	min := strconv.Itoa(int(now.Minute()))
	sec := strconv.Itoa(int(now.Second()))
	myFormat := year + "-" + month + "-" + day + "-" + hour + "-" + min + "-" + sec

	svc := s3.New(session.New())

	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(strings.NewReader(s)),
		Bucket: aws.String(BUCKETNAME),

		// Filename will be in the following format: exampleobject-YYYY-MM-DD-HH-MM-SS
		Key:     aws.String("exampleobject-" + myFormat),
		Tagging: aws.String("env=dev&owner=synapticcleft"),
	}

	result, err := svc.PutObject(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c <- result
	close(c)
}

func sendToDynamo(id string, data string, c chan *dynamodb.PutItemOutput) {
	TABLENAME := os.Getenv("TABLE")
	svc := dynamodb.New(session.New())

	input := &dynamodb.PutItemInput{
		TableName: aws.String(TABLENAME),
		Item: map[string]*dynamodb.AttributeValue{
			"EVENT_ID": {
				S: aws.String(id),
			},
			"DATA": {
				S: aws.String(data),
			},
		},
	}

	result, err := svc.PutItem(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c <- result
	close(c)
}
