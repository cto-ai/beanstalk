package setup

import (
	"fmt"
	"os"

	"git.cto.ai/provision/internal/logger"
	"github.com/aws/aws-sdk-go/aws/session"
	ctoai "github.com/cto-ai/sdk-go"
)

// SDK
type SDKClients struct {
	Ux     *ctoai.Ux
	Prompt *ctoai.Prompt
	Sdk    *ctoai.Sdk
}

// GITHUB

type GithubRepoDetails struct {
	Username string
	Token    string
	Repo     string
	Platform string
}

func PrintIntro(opsClients *SDKClients) error {
	switch opsClients.Sdk.GetInterfaceType() {
	case "terminal":
		logger.LogSlack(opsClients.Ux, `
  [94m██████[39m[33m╗[39m [94m████████[39m[33m╗[39m  [94m██████[39m[33m╗ [39m      [94m█████[39m[33m╗[39m  [94m██[39m[33m╗[39m
 [94m██[39m[33m╔════╝[39m [33m╚══[39m[94m██[39m[33m╔══╝[39m [94m██[39m[33m╔═══[39m[94m██[39m[33m╗[39m     [94m██[39m[33m╔══[39m[94m██[39m[33m╗[39m [94m██[39m[33m║[39m
 [94m██[39m[33m║     [39m [94m   ██[39m[33m║   [39m [94m██[39m[33m║[39m[94m   ██[39m[33m║[39m     [94m███████[39m[33m║[39m [94m██[39m[33m║[39m
 [94m██[39m[33m║     [39m [94m   ██[39m[33m║   [39m [94m██[39m[33m║[39m[94m   ██[39m[33m║[39m     [94m██[39m[33m╔══[39m[94m██[39m[33m║[39m [94m██[39m[33m║[39m
 [33m╚[39m[94m██████[39m[33m╗[39m [94m   ██[39m[33m║   [39m [33m╚[39m[94m██████[39m[33m╔╝[39m [94m██[39m[33m╗[39m [94m██[39m[33m║[39m[94m  ██[39m[33m║[39m [94m██[39m[33m║[39m
 [33m ╚═════╝[39m [33m   ╚═╝   [39m [33m ╚═════╝ [39m [33m╚═╝[39m [33m╚═╝  ╚═╝[39m [33m╚═╝[39m
We’re building the world’s best developer experiences.
  `)
	default:
		logger.LogSlack(opsClients.Ux, `:white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square:
:white_square::white_square::black_square::black_square::white_square::white_square::black_square::black_square::black_square::white_square::white_square::white_square::black_square::black_square::black_square::white_square:
:white_square::black_square::white_square::white_square::black_square::white_square::black_square::white_square::white_square::black_square::white_square::black_square::white_square::white_square::white_square::white_square:
:white_square::black_square::white_square::white_square::black_square::white_square::black_square::black_square::black_square::white_square::white_square::white_square::black_square::black_square::white_square::white_square:
:white_square::black_square::white_square::white_square::black_square::white_square::black_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::black_square::white_square:
:white_square::white_square::black_square::black_square::white_square::white_square::black_square::white_square::white_square::white_square::white_square::black_square::black_square::black_square::white_square::white_square:
:white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square:`)
	}

	logger.LogSlack(opsClients.Ux, "\nCTO.ai Ops - Beanstalk\n")
	logger.LogSlack(opsClients.Ux, `This Op will create an Elastic Beanstalk application and deploy your Github repository.
In addition, it can create a Relational Database Service for your Elastic Beanstalk application.

Requirements:
 - Github
    - Username
    - Access Token (If the repository is private.)
    - Repository Name
	
  - AWS
    - Access Key ID
    - Secret Access Key
	`)

	return nil
}

func GithubSetup(opsClients *SDKClients) (GithubRepoDetails, error) {
	githubRepoDetails, err := promptGithubInfo(opsClients)
	if err != nil {
		return githubRepoDetails, err
	}

	githubRepoAccess := "Public"
	if githubRepoDetails.Token != "public" {
		githubRepoAccess = "Private"
	}

	logger.LogSlack(opsClients.Ux, fmt.Sprintf("ℹ️  Github Information: \n   Username: %s\n   Repo: %s\n   RepoAccess: %s\n   Platform: %s", githubRepoDetails.Username, githubRepoDetails.Repo, githubRepoAccess, githubRepoDetails.Platform))

	confirmGithubInfo, err := opsClients.Prompt.Confirm("GITHUB_BOOL", "Please confirm your Github Information", ctoai.OptConfirmFlag("c"), ctoai.OptConfirmDefault(false))
	if err != nil {
		return githubRepoDetails, err
	}

	if !confirmGithubInfo {
		githubRepoDetails, err = GithubSetup(opsClients)
		if err != nil {
			return githubRepoDetails, err
		}
	}

	return githubRepoDetails, nil
}

func PromptEBAction(prompt *ctoai.Prompt) (string, error) {
	elasticBeanstalkActionOptions := []string{
		"Create New",
		"Update Existing",
	}

	elasticBeanstalkAction, err := prompt.List("EB_OP_OPTION", "Would you like to create a new Elastic Beanstalk Application, or update an existing one?", elasticBeanstalkActionOptions, ctoai.OptListDefaultValue("Create New"), ctoai.OptListFlag("L"), ctoai.OptListAutocomplete(false))
	if err != nil {
		return "", err
	}

	return elasticBeanstalkAction, nil
}

func promptGithubInfo(opsClients *SDKClients) (GithubRepoDetails, error) {
	githubRepoDetails := GithubRepoDetails{}

	githubUser, err := opsClients.Prompt.Input("GITHUB_USER_NAME", "Github Username", ctoai.OptInputAllowEmpty(false))
	if err != nil {
		logger.LogSlackError(opsClients.Ux, err)
		return githubRepoDetails, err
	}
	githubRepoDetails.Username = githubUser

	githubRepo, err := opsClients.Prompt.Input("GITHUB_REPO", "Github Repository", ctoai.OptInputAllowEmpty(false))
	if err != nil {
		logger.LogSlackError(opsClients.Ux, err)
		return githubRepoDetails, err
	}
	githubRepoDetails.Repo = githubRepo

	envPlatformChoices := []string{
		"Node",
		"Go",
	}
	envPlatform, err := opsClients.Prompt.List("EB_ENV_PLATFORM", "Elastic Beanstalk Environment Platform", envPlatformChoices, ctoai.OptListDefaultValue("Node"), ctoai.OptListFlag("L"), ctoai.OptListAutocomplete(false))
	if err != nil {
		logger.LogSlackError(opsClients.Ux, err)
		return githubRepoDetails, err
	}
	githubRepoDetails.Platform = envPlatform

	githubRepoPrivate, err := opsClients.Prompt.Confirm("GITHUB_REPO_PUBLIC", "Is this a private repository?", ctoai.OptConfirmFlag("c"), ctoai.OptConfirmDefault(false))
	if err != nil {
		logger.LogSlackError(opsClients.Ux, err)
		return githubRepoDetails, err
	}

	if githubRepoPrivate {
		githubRepoDetails.Token, err = opsClients.Prompt.Secret("GITHUB_ACCESS_TOKEN", "Github Access Token", ctoai.OptSecretFlag("s"))
		if err != nil {
			logger.LogSlackError(opsClients.Ux, err)
			return githubRepoDetails, err
		}
	} else {
		githubRepoDetails.Token = "public"
	}

	return githubRepoDetails, err
}

// AWS

func PromptAWSInfo(prompt *ctoai.Prompt) (string, error) {
	awsAccessKeyID, err := prompt.Secret("AWS_ACCESS_KEY_ID", "AWS_ACCESS_KEY_ID", ctoai.OptSecretFlag("s"))
	if err != nil {
		return "", err
	}

	awsSecretAccessKey, err := prompt.Secret("AWS_SECRET_ACCESS_KEY", "AWS_SECRET_ACCESS_KEY", ctoai.OptSecretFlag("s"))
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

	awsRegion, err := prompt.List("AWS_REGION", "AWS Region", awsRegionChoices, ctoai.OptListDefaultValue("eu-west-1"), ctoai.OptListFlag("L"), ctoai.OptListAutocomplete(false))
	if err != nil {
		return "", err
	}

	os.Setenv("AWS_ACCESS_KEY_ID", awsAccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", awsSecretAccessKey)
	os.Setenv("AWS_REGION", awsRegion)

	return awsRegion, nil
}

func AWSSetup(prompt *ctoai.Prompt) (*session.Session, string, error) {
	awsRegion, err := PromptAWSInfo(prompt)
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
