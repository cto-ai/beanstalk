package main

import (
	"fmt"

	"git.cto.ai/provision/internal/awseb"
	"git.cto.ai/provision/internal/awsrds"
	"git.cto.ai/provision/internal/awss3"
	"git.cto.ai/provision/internal/awsvpc"
	"git.cto.ai/provision/internal/files"
	"git.cto.ai/provision/internal/logger"
	"git.cto.ai/provision/internal/setup"
	"git.cto.ai/sdk-go/pkg/sdk"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func newApp(opsSDK *sdk.Sdk, awsSess *session.Session, githubRepoDetails setup.GithubRepoDetails, awsRegion string) error {
	rdsDetails, rdsBool, err := awsrds.NewRDSSetup(opsSDK, awsSess)
	if err != nil {
		return err
	}

	unzippedRepo, err := files.EBRepoFileSetup(opsSDK, githubRepoDetails, rdsBool, rdsDetails)
	if err != nil {
		return err
	}

	bucketName, err := awss3.EBS3Setup(opsSDK, awsSess, unzippedRepo, awsRegion)
	if err != nil {
		return err
	}

	envName, appName, err := awseb.NewEBAppSetup(opsSDK, awsSess, bucketName, unzippedRepo, githubRepoDetails.Platform, awsRegion)
	if err != nil {
		return err
	}

	ec2Client := ec2.New(awsSess)
	EBEnvSecurityGroupID, err := awsvpc.DescribeEBEnvSecurityGroupID(ec2Client, envName)
	if err != nil {
		return err
	}

	if rdsBool {
		err = awsvpc.AddEBSGToRDSSG(ec2Client, EBEnvSecurityGroupID, rdsDetails.SecurityGroupID)
		if err != nil {
			return err
		}
	}

	logger.LogSlack(opsSDK, fmt.Sprintf("🌐 Elastic Beanstalk Application: https://%s.console.aws.amazon.com/elasticbeanstalk/home?region=%s#/application/overview?applicationName=%s", awsRegion, awsRegion, appName))
	return nil
}

func updateApp(opsSDK *sdk.Sdk, awsSess *session.Session, githubRepoDetails setup.GithubRepoDetails, awsRegion string) error {
	rdsDetails, rdsBool, err := awsrds.UpdateRDSSetup(opsSDK, awsSess, awsRegion)
	if err != nil {
		return err
	}

	unzippedRepo, err := files.EBRepoFileSetup(opsSDK, githubRepoDetails, rdsBool, rdsDetails)
	if err != nil {
		return err
	}

	bucketName, err := awss3.EBS3Setup(opsSDK, awsSess, unzippedRepo, awsRegion)
	if err != nil {
		return err
	}

	appName, err := awseb.UpdateEBAppSetup(opsSDK, awsSess, bucketName, unzippedRepo, awsRegion)
	if err != nil {
		return err
	}

	logger.LogSlack(opsSDK, fmt.Sprintf("🌐 Elastic Beanstalk Application: https://%s.console.aws.amazon.com/elasticbeanstalk/home?region=%s#/application/overview?applicationName=%s", awsRegion, awsRegion, appName))
	return nil
}

func main() {
	opsSDK := sdk.New()
	logger.LogSlack(opsSDK, "\nCTO.ai Ops - Beanstalk\n")
	logger.LogSlack(opsSDK, `This Op will create an Elastic Beanstalk application and deploy your Github repository.
This Op can also create a Relational Database Service for your Elastic Beanstalk application.

Requirements:
 - Github
    - Username
    - Access Token (If the repository is private.)
    - Repository Name
	
  - AWS
    - Access Key ID
    - Secret Access Key
	`)

	githubRepoDetails, err := setup.GithubSetup(opsSDK)
	if err != nil {
		logger.LogSlackError(opsSDK, err)
		return
	}

	awsSess, awsRegion, err := setup.AWSSetup(opsSDK)
	if err != nil {
		logger.LogSlackError(opsSDK, err)
		return
	}

	elasticBeanstalkAction, err := setup.PromptEBAction(opsSDK)
	if err != nil {
		logger.LogSlackError(opsSDK, err)
		return
	}

	if elasticBeanstalkAction == "Create New" {
		err := newApp(opsSDK, awsSess, githubRepoDetails, awsRegion)
		if err != nil {
			logger.LogSlackError(opsSDK, err)
			return
		}
	} else {
		err := updateApp(opsSDK, awsSess, githubRepoDetails, awsRegion)
		if err != nil {
			logger.LogSlackError(opsSDK, err)
			return
		}
	}
}
