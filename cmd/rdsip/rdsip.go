package rdsip

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// GetRDSInstanceIP returns the IP address of an AWS RDS instance given its name
func GetRDSInstanceIP(instanceName string, region string) (string, error) {
	// Create an AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return "", err
	}

	// Create an RDS client
	rdsClient := rds.NewFromConfig(cfg)

	// Describe DB instances
	params := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(instanceName),
	}

	resp, err := rdsClient.DescribeDBInstances(context.TODO(), params)
	if err != nil {
		return "", err
	}

	// Check if the instance is found
	if len(resp.DBInstances) == 0 {
		return "", fmt.Errorf("RDS instance not found: %s", instanceName)
	}

	// Get the endpoint of the RDS instance
	endpoint := resp.DBInstances[0].Endpoint

	return fmt.Sprintf("%s:%d", *endpoint.Address, *endpoint.Port), nil
}

// RdsIP returns the IP address of an AWS RDS instance given its name
func RdsIP() (string, string) {
	// Replace these values with your own AWS region and RDS instance details
	instanceName := "cts-cts-sandbox"
	region := "us-east-1b"

	// Get the IP address of the RDS instance
	ip, err := GetRDSInstanceIP(instanceName, region)
	if err != nil {
		fmt.Println("Error:", err)
		return "", ""
	}

	return "RDS instance IP:", ip
}
