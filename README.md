![](https://cto.ai/static/oss-banner.png)

# Beanstalk Op (GO)

Dynamically deploy applications to AWS Elastic Beanstalk. This Op also streamlines connections to the Amazon Relational Database Service (RDS) for use with deployed applications.

## Requirements

To run this or any other Op, install the [Ops Platform](https://cto.ai/platform).

Find information about how to run and build Ops via the [Ops Platform Documentation](https://cto.ai/docs/overview).

This Op also requires AWS credentials to work with your account. It also requires the GitHub username and the GitHub repository name of the repository for deployment. It may require a GitHub access token if the repository is private. Here's what you'll need before running this Op the first time:

- **AWS Access Key Id**: via the [AWS Management Console](https://console.aws.amazon.com/):
  - `AWS Management Console` -> `Security Credentials` -> `Access Keys`
- **AWS Access Key Secret**: via the [AWS Management Console](https://console.aws.amazon.com/):
  - `AWS Management Console` -> `Security Credentials` -> `Access Keys`
- **AWS IAM Elastic Beanstalk Permissions** via [AWS Management Console](https://console.aws.amazon.com/):
  - `AWS Management Console` -> `Services` -> `IAM`
- **Github Access Token** [GitHub](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line)
- **Github Username** [GitHub](https://help.github.com/en/github/setting-up-and-managing-your-github-user-account/remembering-your-github-username-or-email)
- **Github Repository Name** [GitHub](https://help.github.com/en/github/getting-started-with-github/create-a-repo)

This Op can create and connect RDS database instances to your application. If this is desired, the user will need to provide or create the following information:

- **RDS Database Instance Name**
- **RDS Database Master Username**
- **RDS Database Master Password**

When connecting a RDS database to your application, this Op will create a directory and a file containing your database access information (`.ebextensions/rds_env`) within your application before the deployment. This step can be skipped, however you may be required to connect your application to the RDS instance on your own.

## Usage

To start this Op prompt run:

```bash
ops run @cto.ai/beanstalk
```

## Demo Applications

Example applications that can be deployed with this Op:

- [Node JS Demo Application](https://github.com/eddingston/ops-beanstalk-node-demo)

## Local Development / Running from Source

**1. Clone the repo:**

```bash
git clone <git url>
```

**2. Install dependencies:**

```bash
go get -u ./...
```

**3. Run the Op from your current working directory with:**

```bash
ops run .
```

### AWS Docs

- [Getting Started on Amazon Web Services (AWS)](https://aws.amazon.com/getting-started/)
- [AWS SDK for Go API Reference](https://docs.aws.amazon.com/sdk-for-go/api/)

## Contributors

<table>
  <tr>
    <td align="center"><a href="https://github.com/eddingston"><img src="https://avatars0.githubusercontent.com/u/40420154?s=460&v=4" width="100px;" alt=""/><br /><sub><b>Edmond Lee</b></sub></a><br/></td>
  </tr>
</table>

## LICENSE

[MIT](LICENSE.txt)
