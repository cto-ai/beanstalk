package awseb

import (
	"fmt"
	"strings"
	"time"

	"git.cto.ai/provision/internal/logger"
	"git.cto.ai/sdk-go/pkg/sdk"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
)

func NewEBAppSetup(sdk *sdk.Sdk, awsSess *session.Session, bucketName, unzippedRepo, repoPlatform, awsRegion string) (string, string, error) {
	ebClient := elasticbeanstalk.New(awsSess, aws.NewConfig().WithRegion(awsRegion))

	EBAppName, err := createApp(sdk, ebClient, unzippedRepo)
	if err != nil {
		return "", EBAppName, err
	}

	envName, err := createEnviro(sdk, ebClient, bucketName, EBAppName, repoPlatform)
	if err != nil {
		return envName, EBAppName, err
	}

	err = createAppVersion(sdk, ebClient, EBAppName, bucketName, unzippedRepo)
	if err != nil {
		return envName, EBAppName, err
	}

	err = updateEnvironment(sdk, ebClient, bucketName, envName, 0)
	if err != nil {
		return envName, EBAppName, err
	}

	logger.LogSlack(sdk, fmt.Sprintf("ℹ️  EB Application Name: %s\nℹ️  EB Environment Name: %s\nℹ️  EB Application Version Name: %s", EBAppName, envName, bucketName))

	return envName, EBAppName, nil
}

func PromptEBInfo(sdk *sdk.Sdk, ebClient *elasticbeanstalk.ElasticBeanstalk) (string, string, error) {
	EBAppNameMatches, err := GetSpecifiedEBApps(ebClient)
	if err != nil {
		return "", "", err
	}

	EBAppName, err := sdk.PromptList(EBAppNameMatches, "EB_APP_NAME", "Choose the Elastic Beanstalk app that you want to update, or enter the name of the app", "Enter a value", "flag")
	if err != nil {
		return EBAppName, "", err
	}

	EBAppEnvMatches, err := getSpecifiedEBAppEnv(ebClient, EBAppName)
	if err != nil {
		return EBAppName, "", err
	}

	EBAppEnvName, err := sdk.PromptList(EBAppEnvMatches, "EB_APP_NAME", "Choose the Elastic Beanstalk app environment that you want to update, or enter the name of the environment", "Enter a value", "flag")
	if err != nil {
		return EBAppName, EBAppEnvName, err
	}

	return EBAppName, EBAppEnvName, nil
}

func UpdateEBAppSetup(sdk *sdk.Sdk, awsSess *session.Session, bucketName, unzippedRepo, awsRegion string) (string, error) {
	ebClient := elasticbeanstalk.New(awsSess, aws.NewConfig().WithRegion(awsRegion))

	EBAppName, EBAppEnvName, err := PromptEBInfo(sdk, ebClient)
	if err != nil {
		return EBAppName, err
	}

	err = createAppVersion(sdk, ebClient, EBAppName, bucketName, unzippedRepo)
	if err != nil {
		return EBAppName, err
	}

	err = updateEnvironment(sdk, ebClient, bucketName, EBAppEnvName, 0)
	if err != nil {
		return EBAppName, err
	}

	logger.LogSlack(sdk, fmt.Sprintf("ℹ️  EB Application Name: %s\nℹ️  EB Environment Name: %s\nℹ️  EB Application Version Name: %s", EBAppName, EBAppEnvName, bucketName))

	return EBAppName, nil
}

func GetSpecifiedEBApps(ebClient *elasticbeanstalk.ElasticBeanstalk) ([]string, error) {
	EBAppNameMatches := []string{"Enter a value"}

	result, err := ebClient.DescribeApplications(&elasticbeanstalk.DescribeApplicationsInput{})
	if err != nil {
		return EBAppNameMatches, err
	}

	for _, k := range result.Applications {
		EBAppNameMatches = append(EBAppNameMatches, *k.ApplicationName)
	}

	return EBAppNameMatches, nil
}

func getSpecifiedEBAppEnv(ebClient *elasticbeanstalk.ElasticBeanstalk, ebAppName string) ([]string, error) {
	EBEnvNameMatches := []string{"Enter a value"}

	input := &elasticbeanstalk.DescribeEnvironmentsInput{
		ApplicationName: aws.String(ebAppName),
	}

	result, err := ebClient.DescribeEnvironments(input)
	if err != nil {
		return EBEnvNameMatches, err
	}

	for _, k := range result.Environments {
		EBEnvNameMatches = append(EBEnvNameMatches, *k.EnvironmentName)
	}

	return EBEnvNameMatches, nil
}

func createApp(sdk *sdk.Sdk, ebClient *elasticbeanstalk.ElasticBeanstalk, unzippedRepo string) (string, error) {
	logger.LogSlack(sdk, "🔄 Creating Elastic Beanstalk application...")

	unzippedRepoSplit := strings.Split(unzippedRepo, "-")
	EBAppName := strings.Join(unzippedRepoSplit[:(len(unzippedRepoSplit)-1)], "-")

	input := &elasticbeanstalk.CreateApplicationInput{
		ApplicationName: aws.String(EBAppName),
		Description:     aws.String(EBAppName),
	}

	_, err := ebClient.CreateApplication(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Message() == fmt.Sprintf("Application %s already exists.", EBAppName) {
				logger.LogSlack(sdk, fmt.Sprintf("ℹ️  Application %s already exists. \nℹ️  Skipping to next step...", EBAppName))
				return EBAppName, nil
			} else {
				return EBAppName, aerr
			}
		} else {
			return EBAppName, err
		}
	}

	logger.LogSlack(sdk, "✅ Elastic Beanstalk application created.")
	return EBAppName, nil
}

func createEnviro(sdk *sdk.Sdk, ebClient *elasticbeanstalk.ElasticBeanstalk, bucketName, EBAppName, envPlatform string) (string, error) {
	logger.LogSlack(sdk, "🔄 Creating Elastic Beanstalk application environment...")

	bucketNameSplit := strings.Split(bucketName, "-")
	envName := bucketNameSplit[len(bucketNameSplit)-1]

	input := &elasticbeanstalk.CreateEnvironmentInput{
		ApplicationName:   aws.String(EBAppName),
		CNAMEPrefix:       aws.String(bucketName),
		EnvironmentName:   aws.String(envName),
		SolutionStackName: aws.String("64bit Amazon Linux 2018.03 v2.14.2 running Go 1.13.6"),
	}

	switch envPlatform {
	case "Go":
		input.SolutionStackName = aws.String("64bit Amazon Linux 2018.03 v2.14.2 running Go 1.13.6")

	case "Node":
		configOptionSetting := elasticbeanstalk.ConfigurationOptionSetting{
			Namespace:  aws.String("aws:elasticbeanstalk:container:nodejs"),
			OptionName: aws.String("NodeCommand"),
			Value:      aws.String("npm start"),
		}
		input.SolutionStackName = aws.String("64bit Amazon Linux 2018.03 v4.13.0 running Node.js")
		input.OptionSettings = []*elasticbeanstalk.ConfigurationOptionSetting{&configOptionSetting}
	}

	_, err := ebClient.CreateEnvironment(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return envName, aerr
		} else {
			return "", err
		}
	}

	logger.LogSlack(sdk, "✅ Elastic Beanstalk application environment created.")
	return envName, nil
}

func createAppVersion(sdk *sdk.Sdk, ebClient *elasticbeanstalk.ElasticBeanstalk, EBAppName, bucketName, unzippedRepo string) error {
	logger.LogSlack(sdk, "🔄 Creating Elastic Beanstalk application version...")

	input := &elasticbeanstalk.CreateApplicationVersionInput{
		ApplicationName:       aws.String(EBAppName),
		AutoCreateApplication: aws.Bool(true),
		Description:           aws.String(bucketName),
		Process:               aws.Bool(true),
		SourceBundle: &elasticbeanstalk.S3Location{
			S3Bucket: aws.String(bucketName),
			S3Key:    aws.String(fmt.Sprintf("%s.zip", unzippedRepo)),
		},
		VersionLabel: aws.String(bucketName),
	}

	_, err := ebClient.CreateApplicationVersion(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr
		}
		return err
	}

	logger.LogSlack(sdk, "✅ Elastic Beanstalk application version created.")
	return nil
}

func updateEnvironment(sdk *sdk.Sdk, svc *elasticbeanstalk.ElasticBeanstalk, bucketName, envName string, retries int) error {
	if retries%2 == 0 {
		logger.LogSlack(sdk, "🔄 Preparing to update Elastic Beanstalk application environment...")
	}

	input := &elasticbeanstalk.UpdateEnvironmentInput{
		EnvironmentName: aws.String(envName),
		VersionLabel:    aws.String(bucketName),
	}
	_, err := svc.UpdateEnvironment(input)
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == "InvalidParameterValue" && aerr.Message() == fmt.Sprintf("Environment named %s is in an invalid state for this operation. Must be Ready.", envName) && retries <= 20 {
			time.Sleep(30 * time.Second)

			err := updateEnvironment(sdk, svc, bucketName, envName, retries+1)
			if err != nil {
				return aerr
			}
		} else {
			return aerr
		}
	} else {
		return err
	}

	if retries == 0 {
		logger.LogSlack(sdk, "✅ Elastic Beanstalk application environment has started to update, please wait for it to finish.")

	}

	return nil
}
