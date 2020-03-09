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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	ctoai "github.com/cto-ai/sdk-go"
)

func newApp(opsClients *setup.SDKClients, awsSess *session.Session, githubRepoDetails setup.GithubRepoDetails, awsRegion string) error {
	rdsDetails, rdsBool, err := awsrds.NewRDSSetup(opsClients, awsSess)
	if err != nil {
		return err
	}

	unzippedRepo, err := files.EBRepoFileSetup(opsClients.Ux, githubRepoDetails, rdsBool, rdsDetails)
	if err != nil {
		return err
	}

	bucketName, err := awss3.EBS3Setup(opsClients.Ux, awsSess, unzippedRepo, awsRegion)
	if err != nil {
		return err
	}

	envName, appName, err := awseb.NewEBAppSetup(opsClients.Ux, awsSess, bucketName, unzippedRepo, githubRepoDetails.Platform, awsRegion)
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

	logger.LogSlack(opsClients.Ux, fmt.Sprintf("🌐 Elastic Beanstalk Application: https://%s.console.aws.amazon.com/elasticbeanstalk/home?region=%s#/application/overview?applicationName=%s", awsRegion, awsRegion, appName))
	return nil
}

func updateApp(opsClients *setup.SDKClients, awsSess *session.Session, githubRepoDetails setup.GithubRepoDetails, awsRegion string) error {
	rdsDetails, rdsBool, err := awsrds.UpdateRDSSetup(opsClients, awsSess, awsRegion)
	if err != nil {
		return err
	}

	unzippedRepo, err := files.EBRepoFileSetup(opsClients.Ux, githubRepoDetails, rdsBool, rdsDetails)
	if err != nil {
		return err
	}

	bucketName, err := awss3.EBS3Setup(opsClients.Ux, awsSess, unzippedRepo, awsRegion)
	if err != nil {
		return err
	}

	appName, err := awseb.UpdateEBAppSetup(opsClients, awsSess, bucketName, unzippedRepo, awsRegion)
	if err != nil {
		return err
	}

	logger.LogSlack(opsClients.Ux, fmt.Sprintf("🌐 Elastic Beanstalk Application: https://%s.console.aws.amazon.com/elasticbeanstalk/home?region=%s#/application/overview?applicationName=%s", awsRegion, awsRegion, appName))
	return nil
}

func main() {
	opsClients := setup.SDKClients{
		Ux:     ctoai.NewUx(),
		Prompt: ctoai.NewPrompt(),
		Sdk:    ctoai.NewSdk(),
	}

	err := setup.PrintIntro(&opsClients)
	if err != nil {
		logger.LogSlackError(opsClients.Ux, err)
		return
	}

	githubRepoDetails, err := setup.GithubSetup(&opsClients)
	if err != nil {
		logger.LogSlackError(opsClients.Ux, err)
		return
	}

	awsSess, awsRegion, err := setup.AWSSetup(opsClients.Prompt)
	if err != nil {
		logger.LogSlackError(opsClients.Ux, err)
		return
	}

	elasticBeanstalkAction, err := setup.PromptEBAction(opsClients.Prompt)
	if err != nil {
		logger.LogSlackError(opsClients.Ux, err)
		return
	}

	if elasticBeanstalkAction == "Create New" {
		err := newApp(&opsClients, awsSess, githubRepoDetails, awsRegion)
		if err != nil {
			logger.LogSlackError(opsClients.Ux, err)
			return
		}
	} else {
		err := updateApp(&opsClients, awsSess, githubRepoDetails, awsRegion)
		if err != nil {
			logger.LogSlackError(opsClients.Ux, err)
			return
		}
	}
}
