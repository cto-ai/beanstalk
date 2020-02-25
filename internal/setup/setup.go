package setup

import (
	"fmt"
	"os"

	"git.cto.ai/provision/internal/logger"
	"git.cto.ai/sdk-go/pkg/sdk"
	"github.com/aws/aws-sdk-go/aws/session"
)

// GITHUB

type GithubRepoDetails struct {
	Username string
	Token    string
	Repo     string
	Platform string
}

func GithubSetup(sdk *sdk.Sdk) (GithubRepoDetails, error) {
	githubRepoDetails, err := promptGithubInfo(sdk)
	if err != nil {
		return githubRepoDetails, err
	}

	githubRepoAccess := "Public"
	if githubRepoDetails.Token != "public" {
		githubRepoAccess = "Private"
	}

	logger.LogSlack(sdk, fmt.Sprintf("ℹ️  Github Information: \n   Username: %s\n   Repo: %s\n   RepoAccess: %s\n   Platform: %s", githubRepoDetails.Username, githubRepoDetails.Repo, githubRepoAccess, githubRepoDetails.Platform))

	confirmGithubInfo, err := sdk.PromptConfirm("GITHUB_BOOL", "Please confirm your Github Information", "flag", false)
	if err != nil {
		return githubRepoDetails, err
	}

	if !confirmGithubInfo {
		githubRepoDetails, err = GithubSetup(sdk)
		if err != nil {
			return githubRepoDetails, err
		}
	}

	return githubRepoDetails, nil
}

func PromptEBAction(sdk *sdk.Sdk) (string, error) {
	elasticBeanstalkActionOptions := []string{
		"Create New",
		"Update Existing",
	}

	elasticBeanstalkAction, err := sdk.PromptList(elasticBeanstalkActionOptions, "EB_OP_OPTION", "Would you like to create a new Elastic Beanstalk Application, or update an existing one?", "Create New", "flag")
	if err != nil {
		return "", err
	}

	return elasticBeanstalkAction, nil
}

func promptGithubInfo(sdk *sdk.Sdk) (GithubRepoDetails, error) {
	githubRepoDetails := GithubRepoDetails{}

	githubUser, err := sdk.PromptInput("GITHUB_USER_NAME", "Github Username", "", "flag", false)
	if err != nil {
		logger.LogSlackError(sdk, err)
		return githubRepoDetails, err
	}
	githubRepoDetails.Username = githubUser

	githubRepo, err := sdk.PromptInput("GITHUB_REPO", "Github Repository", "", "flag", false)
	if err != nil {
		logger.LogSlackError(sdk, err)
		return githubRepoDetails, err
	}
	githubRepoDetails.Repo = githubRepo

	envPlatformChoices := []string{
		"Node",
		"Go",
	}
	envPlatform, err := sdk.PromptList(envPlatformChoices, "EB_ENV_PLATFORM", "Elastic Beanstalk Environment Platform", "Node", "flag")
	if err != nil {
		logger.LogSlackError(sdk, err)
		return githubRepoDetails, err
	}
	githubRepoDetails.Platform = envPlatform

	githubRepoPrivate, err := sdk.PromptConfirm("GITHUB_REPO_PUBLIC", "Is this a private repository?", "flag", true)
	if err != nil {
		logger.LogSlackError(sdk, err)
		return githubRepoDetails, err
	}

	if githubRepoPrivate {
		githubRepoDetails.Token, err = sdk.PromptSecret("GITHUB_ACCESS_TOKEN", "Github Access Token", "flag")
		if err != nil {
			logger.LogSlackError(sdk, err)
			return githubRepoDetails, err
		}
	} else {
		githubRepoDetails.Token = "public"
	}

	return githubRepoDetails, err
}

// AWS

func PromptAWSInfo(sdk *sdk.Sdk) (string, error) {
	awsAccessKeyID, err := sdk.PromptSecret("AWS_ACCESS_KEY_ID", "AWS_ACCESS_KEY_ID", "flag")
	if err != nil {
		return "", err
	}

	awsSecretAccessKey, err := sdk.PromptSecret("AWS_SECRET_ACCESS_KEY", "AWS_SECRET_ACCESS_KEY", "flag")
	if err != nil {
		return "", err
	}
	awsRegionChoices := []string{"us-east-2",
		"us-east-1",
		"us-west-1",
		"us-west-2",
		"ap-east-1",
		"ap-south-1",
		"ap-northeast-3",
		"ap-northeast-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-northeast-1",
		"ca-central-1",
		"cn-north-1",
		"cn-northwest-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"eu-north-1",
		"me-south-1",
		"sa-east-1",
		"us-gov-east-1",
		"us-gov-west-1"}

	awsRegion, err := sdk.PromptList(awsRegionChoices, "AWS_REGION", "AWS Region", "eu-west-1", "flag")
	if err != nil {
		return "", err
	}

	os.Setenv("AWS_ACCESS_KEY_ID", awsAccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", awsSecretAccessKey)
	os.Setenv("AWS_REGION", awsRegion)

	return awsRegion, nil
}

func AWSSetup(sdk *sdk.Sdk) (*session.Session, string, error) {
	awsRegion, err := PromptAWSInfo(sdk)
	if err != nil {
		return nil, "", err
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, "", err
	}

	return sess, awsRegion, nil
}
