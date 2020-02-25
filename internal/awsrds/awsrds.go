package awsrds

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"

	"git.cto.ai/provision/internal/logger"
	"git.cto.ai/sdk-go/pkg/sdk"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSDetails struct {
	Username        string
	Password        string
	Host            string
	Port            string
	DBName          string
	SubnetGroup     string
	SecurityGroupID string
	Platform        string
}

func confirmRDSPassword(sdk *sdk.Sdk, statement string) (string, error) {
	var confirmedPassword string

	password, err := sdk.PromptSecret("RDS_DB_PASSWORD", statement, "flag")
	if err != nil {
		return confirmedPassword, err
	}

	confirmPassword, err := sdk.PromptSecret("RDS_DB_PASSWORD", statement, "flag")
	if err != nil {
		return confirmedPassword, err
	}

	if password != confirmPassword {
		logger.LogSlack(sdk, fmt.Sprintf("ℹ️  The passwords did not match. Please try again."))
		confirmedPassword, err = confirmRDSPassword(sdk, statement)
		if err != nil {
			return confirmPassword, err
		}
	} else {
		confirmedPassword = password
	}

	return confirmedPassword, nil
}

func NewRDSSetup(sdk *sdk.Sdk, awsSess *session.Session) (RDSDetails, bool, error) {
	rdsDetails, rdsBool, err := setRDSInfo(sdk)
	if err != nil {
		return rdsDetails, rdsBool, err
	}

	logger.LogSlack(sdk, fmt.Sprintf("ℹ️  RDS Information: \n   DBName: %s\n   MasterUsername: %s\n   Platform: %s", rdsDetails.DBName, rdsDetails.Username, rdsDetails.Platform))

	confirmRDSInfo, err := sdk.PromptConfirm("RDS_BOOL", "Please confirm your RDS Information", "flag", false)
	if err != nil {
		return rdsDetails, rdsBool, err
	}

	if !confirmRDSInfo {
		rdsDetails, rdsBool, err = NewRDSSetup(sdk, awsSess)
		if err != nil {
			return rdsDetails, rdsBool, err
		}
	}

	rdsClient := rds.New(awsSess)

	RDSSecurityGroupID, err := createRDSInstance(sdk, rdsClient, rdsDetails)
	if err != nil {
		return rdsDetails, rdsBool, err
	}
	rdsDetails.SecurityGroupID = RDSSecurityGroupID

	dbHost, dbPort, err := getSpecifiedDBInstanceEndpoint(sdk, rdsClient, rdsDetails.DBName, 0)
	rdsDetails.Host = dbHost
	rdsDetails.Port = dbPort

	return rdsDetails, rdsBool, nil
}

func setRDSInfo(sdk *sdk.Sdk) (RDSDetails, bool, error) {
	rdsDetails := RDSDetails{}

	rdsBool, err := sdk.PromptConfirm("RDS_BOOL", "Does your app require a RDS database instance?", "flag", false)
	if err != nil {
		return rdsDetails, rdsBool, err
	}

	if rdsBool {
		rdsBool, err = sdk.PromptConfirm("RDS_BOOL", "A directory and a file containing the RDS details variables will be added to your repository in order to connect to the RDS database instance. If this is acceptable, press 'Y' to continue, otherwise you may continue without an RDS database instance.", "flag", false)
		if err != nil {
			return rdsDetails, rdsBool, err
		}
	}

	if rdsBool {
		dbIdentifierID, err := sdk.PromptInput("DB_INSTANCE_IDENTIFIER", "RDS Instance Name", "", "flag", false)
		if err != nil {
			return rdsDetails, rdsBool, err
		}
		rdsDetails.DBName = dbIdentifierID

		rdsPlatformChoices := []string{
			"postgres",
		}
		rdsPlatform, err := sdk.PromptList(rdsPlatformChoices, "RDS_PLATFORM", "RDS Platform", "postgres", "flag")
		if err != nil {
			return rdsDetails, rdsBool, err
		}
		rdsDetails.Platform = rdsPlatform

		dbUsername, err := sdk.PromptInput("RDS_DB_USERNAME", "RDS DB Username", "", "flag", false)
		if err != nil {
			return rdsDetails, rdsBool, err
		}
		rdsDetails.Username = dbUsername

		dbPassword, err := confirmRDSPassword(sdk, "RDS DB Password")
		if err != nil {
			return rdsDetails, rdsBool, err
		}
		rdsDetails.Password = dbPassword
	}

	return rdsDetails, rdsBool, nil
}

func UpdateRDSSetup(sdk *sdk.Sdk, awsSess *session.Session, awsRegion string) (RDSDetails, bool, error) {
	rdsDetails, rdsBool, err := getRDSInfo(sdk, awsSess, awsRegion)
	if err != nil {
		return rdsDetails, rdsBool, err
	}

	logger.LogSlack(sdk, fmt.Sprintf("ℹ️  RDS Information: \n   DBName: %s\n   MasterUsername: %s\n   Host: %s\n   Port: %s", rdsDetails.DBName, rdsDetails.Username, rdsDetails.Host, rdsDetails.Port))

	confirmRDSInfo, err := sdk.PromptConfirm("RDS_BOOL", "Please confirm your RDS Information", "flag", false)
	if err != nil {
		return rdsDetails, rdsBool, err
	}

	if !confirmRDSInfo {
		rdsDetails, rdsBool, err = UpdateRDSSetup(sdk, awsSess, awsRegion)
		if err != nil {
			return rdsDetails, rdsBool, err
		}
	}

	return rdsDetails, rdsBool, nil
}

func getRDSInfo(sdk *sdk.Sdk, awsSess *session.Session, awsRegion string) (RDSDetails, bool, error) {
	rdsBool, err := sdk.PromptConfirm("RDS_BOOL", "Does your app require a RDS database?", "flag", false)
	if err != nil {
		return RDSDetails{}, rdsBool, err
	}

	rdsDetails := RDSDetails{}

	var rdsExisting bool
	rdsClient := rds.New(awsSess)

	if rdsBool {
		rdsExisting, err = sdk.PromptConfirm("RDS_BOOL", "Does your already have an existing RDS database?", "flag", false)
		if err != nil {
			return RDSDetails{}, rdsBool, err
		}
	}

	if rdsExisting {
		rdsInstanceNameMatches, rdsInstances, err := getAllRDSInstanceNames(sdk, rdsClient)

		rdsInstanceName, err := sdk.PromptList(rdsInstanceNameMatches, "RDS_INSTANCE_NAME", "Please choose the RDS instance that is connected to the app", "Enter a value", "flag")
		if err != nil {
			return RDSDetails{}, rdsBool, err
		}

		for _, k := range rdsInstances {
			rdsInstanceNameMatches = append(rdsInstanceNameMatches, *k.DBInstanceIdentifier)

			if rdsInstanceName == *k.DBInstanceIdentifier {
				rdsDetails.DBName = *k.DBInstanceIdentifier
				rdsDetails.Username = *k.MasterUsername
				rdsDetails.Host = *k.Endpoint.Address
				rdsDetails.Port = fmt.Sprintf("%v", *k.Endpoint.Port)
			}
			continue
		}

		dbPassword, err := confirmRDSPassword(sdk, "Please enter the master password for the chosen RDS instance")
		if err != nil {
			return rdsDetails, rdsBool, err
		}
		rdsDetails.Password = dbPassword
	}

	return rdsDetails, rdsBool, err
}

func createRDSInstance(sdk *sdk.Sdk, rdsClient *rds.RDS, rdsDetails RDSDetails) (string, error) {
	logger.LogSlack(sdk, "🔄 Creating RDS database...")

	input := &rds.CreateDBInstanceInput{
		AllocatedStorage:     aws.Int64(5),
		DBInstanceClass:      aws.String("db.t2.micro"),
		DBInstanceIdentifier: aws.String(rdsDetails.DBName),
		Engine:               aws.String(rdsDetails.Platform),
		MasterUserPassword:   aws.String(rdsDetails.Password),
		MasterUsername:       aws.String(rdsDetails.Username),
	}

	results, err := rdsClient.CreateDBInstance(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return "", aerr
		}
		return "", err
	}

	logger.LogSlack(sdk, "✅ RDS database created.")
	logger.LogSlack(sdk, "ℹ️  Beginning to set up RDS database. This may take around 5 minutes.")
	logger.LogSlack(sdk, "🔄 Setting up RDS database...")

	return *results.DBInstance.VpcSecurityGroups[0].VpcSecurityGroupId, nil
}

func getSpecifiedDBInstanceEndpoint(sdk *sdk.Sdk, rdsClient *rds.RDS, DBIdentifierID string, retries int64) (string, string, error) {
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(DBIdentifierID),
	}

	result, err := rdsClient.DescribeDBInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return "", "", aerr
		}
		return "", "", err
	}

	var dbHost string
	var dbPort string

	if *result.DBInstances[0].DBInstanceStatus != "available" {
		time.Sleep(15 * time.Second)

		if retries%4 == 0 {
			logger.LogSlack(sdk, "🔄 Setting up RDS database...")
		}

		retryHost, retryPort, err := getSpecifiedDBInstanceEndpoint(sdk, rdsClient, DBIdentifierID, retries+1)
		if err != nil {
			return "", "", err
		}

		dbHost = retryHost
		dbPort = retryPort

	} else {
		logger.LogSlack(sdk, "✅ RDS database setup completed.")

		host := fmt.Sprintf("%v", *result.DBInstances[0].Endpoint.Address)
		port := fmt.Sprintf("%v", *result.DBInstances[0].Endpoint.Port)

		return host, port, nil
	}

	return dbHost, dbPort, err
}

func getAllRDSInstanceNames(sdk *sdk.Sdk, rdsClient *rds.RDS) ([]string, []*rds.DBInstance, error) {
	input := &rds.DescribeDBInstancesInput{}

	rdsInstanceNameMatches := []string{"Enter a value"}

	result, err := rdsClient.DescribeDBInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return rdsInstanceNameMatches, result.DBInstances, aerr
		}
		return rdsInstanceNameMatches, result.DBInstances, err
	}

	for _, k := range result.DBInstances {
		rdsInstanceNameMatches = append(rdsInstanceNameMatches, *k.DBInstanceIdentifier)
	}

	return rdsInstanceNameMatches, result.DBInstances, nil
}
