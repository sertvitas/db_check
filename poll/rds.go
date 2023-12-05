package poll

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"

	// Import the postgres driver
	_ "github.com/lib/pq"
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
func InstanceIsAvailable(instanceID string, log *zerolog.Logger) bool {
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
	log.Info().Msgf("Instance %s status: %s", instanceID, *output.DBInstances[0].DBInstanceStatus)
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
		if InstanceIsAvailable(secret.DbInstanceIdentifier, log) {
			log.Info().Msgf("%s is available", secret.DbInstanceIdentifier)
			*isAvailable = true
			break
		}
		log.Info().Msgf("%s is not available", secret.DbInstanceIdentifier)
		time.Sleep(time.Duration(pollDelaySeconds) * time.Second)
	}
}

// CheckDBConnection attempts to connect to the database using the given credentials
func CheckDBConnection(secret RDSSecretData) error {
	// Create a connection string
	connStr := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=postgres sslmode=disable",
		secret.Host, secret.Port, secret.Username, secret.Password)

	// Attempt to connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// Ping the database to verify connectivity
	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

// CheckTCPConnectivity checks TCP connectivity to the specified host and port
func CheckTCPConnectivity(secret RDSSecretData) error {
	address := fmt.Sprintf("%s:%d", secret.Host, secret.Port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", address, err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	return nil
}

// DBLogin checks TCP connectivity and DB connectivity to the specified host and port
func DBLogin(secret RDSSecretData) (err error) {
	err = CheckTCPConnectivity(secret)
	if err != nil {
		return err
	}
	err = CheckDBConnection(secret)
	return nil
}
