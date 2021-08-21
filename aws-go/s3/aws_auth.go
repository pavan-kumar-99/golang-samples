package main

import (
	"fmt"

	aws "github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	s3 "github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		fmt.Errorf(err.Error())
	}
	s3c := s3.New(sess)
	input := &s3.ListBucketsInput{}
	buckets, err := s3c.ListBuckets(input)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	fmt.Println(buckets)
}
