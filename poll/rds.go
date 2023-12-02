package poll

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/rs/zerolog"
)

// RDSSecretData is the struct of the secret generated for RDS by CDK deployment
// Password: the password for the database user
// Engine: the database engine
// Port: the port the database is listening on
// DbInstanceIdentifier: the unique name of the RDS instance
// Host: the hostname of the RDS instance
// Username: the username for the database user
type RDSSecretData struct {
	Password             string `json:"password"`
	Engine               string `json:"engine"`
	Port                 int    `json:"port"`
	DbInstanceIdentifier string `json:"dbInstanceIdentifier"`
	Host                 string `json:"host"`
	Username             string `json:"username"`
}

// InstanceIsAvailable checks if the RDS instance is available
func InstanceIsAvailable(instanceID string) bool {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	// Create an RDS client
	client := rds.NewFromConfig(cfg)

	// Create a DescribeDBInstancesInput object
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: &instanceID,
	}

	// Retrieve information about the specified RDS instance
	output, err := client.DescribeDBInstances(context.Background(), input)
	if err != nil {
		panic(err)
	}

	// Check if any DBInstances were returned
	if len(output.DBInstances) == 0 {
		panic(fmt.Sprintf("no instances found with ID: %s", instanceID))
	}

	// Extract the status of the RDS instance
	return *output.DBInstances[0].DBInstanceStatus == "available"
}

// InstanceMonitor polls the RDS instance every 60 seconds
func InstanceMonitor(
	isAvailable *bool,
	initialDelaySeconds int,
	pollDelaySeconds int,
	secret RDSSecretData,
	log *zerolog.Logger) {
	log.Info().Msgf(
		"Waiting for %v seconds before monitoring %s",
		initialDelaySeconds, secret.DbInstanceIdentifier)
	time.Sleep(time.Duration(initialDelaySeconds) * time.Second)
	log.Info().Msgf("Monitoring %s", secret.DbInstanceIdentifier)
	for {
		if InstanceIsAvailable(secret.DbInstanceIdentifier) {
			log.Info().Msgf("%s is available", secret.DbInstanceIdentifier)
			*isAvailable = true
			break
		}
		log.Info().Msgf("%s is not available", secret.DbInstanceIdentifier)
		time.Sleep(time.Duration(pollDelaySeconds) * time.Second)
	}
}
