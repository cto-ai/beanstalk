package awsvpc

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func AddEBSGToRDSSG(ec2Client *ec2.EC2, ebSG, rdsSG string) error {
	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(rdsSG),
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort:   aws.Int64(5432),
				IpProtocol: aws.String("tcp"),
				ToPort:     aws.Int64(5432),
				UserIdGroupPairs: []*ec2.UserIdGroupPair{
					{
						GroupId: aws.String(ebSG),
					},
				},
			},
		},
	}

	_, err := ec2Client.AuthorizeSecurityGroupIngress(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr
		}
		return err
	}

	return nil
}

func DescribeEBEnvSecurityGroupID(ec2Client *ec2.EC2, envName string) (string, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:elasticbeanstalk:environment-name"),
				Values: []*string{
					aws.String(envName),
				},
			},
		},
	}

	result, err := ec2Client.DescribeSecurityGroups(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return "", aerr
		}
		return "", err
	}

	var envSecurityGroupID string

	for _, k := range result.SecurityGroups {
		if *k.Description == "SecurityGroup for ElasticBeanstalk environment." {
			envSecurityGroupID = *k.GroupId
			continue
		}
	}

	return envSecurityGroupID, nil
}
