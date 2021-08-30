package main

import (
	"fmt"

	aws "github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	ec2 "github.com/aws/aws-sdk-go/service/ec2"
)

var sess *ec2.EC2

func init() {
	sess = ec2.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
	})))
}

func main() {
	volid := "vol-00e39e3eb62630e69"
	input := &ec2.DescribeVolumesInput{
		VolumeIds: []*string{&volid},
	}
	vol, err := sess.DescribeVolumes(input)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(*vol.Volumes[0].AvailabilityZone)
}
