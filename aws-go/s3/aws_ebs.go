package aws

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
	fmt.Println(filter_vol_by_id(volid))
	fmt.Println(filter_vol_by_tags("Name", "test"))
}

func filter_vol_by_id(volid string) (string, error) {

	input := &ec2.DescribeVolumesInput{
		VolumeIds: []*string{&volid},
	}
	vol, err := sess.DescribeVolumes(input)
	if err != nil {
		return "", err
	}
	vol_az := *vol.Volumes[0].AvailabilityZone
	return vol_az, nil
}

func filter_vol_by_tags(tag_key, tag_value string) (string, error) {
	input := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(tag_value)},
			},
		},
	}
	vol, err := sess.DescribeVolumes(input)
	if err != nil {
		return "", err
	}
	vol_az := *&vol.Volumes[0].VolumeId
	return *vol_az, nil
}

//To-Do add Caching mechanism to cache the data from AWS
//To add filtering of EBS volumes from multiple tags
