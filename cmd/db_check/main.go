package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/sertvitas/db_check/report"
	"net"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/sertvitas/db_check/version"

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

// CheckTCPConnectivity checks TCP connectivity to the specified host and port
func CheckTCPConnectivity(secret RDSSecretData) error {
	address := fmt.Sprintf("%s:%d", secret.Host, secret.Port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", address, err)
	}
	defer conn.Close()
	return nil
}

// GetRDSSecret retrieves the secret value for the given secret ID and unmarshals it into RDSSecretData
func GetRDSSecret(secretID string) (*RDSSecretData, error) {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	// Create a Secrets Manager client
	client := secretsmanager.NewFromConfig(cfg)

	// Input for the GetSecretValue API call
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretID,
	}

	// Execute the GetSecretValue API call
	result, err := client.GetSecretValue(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret value, %v", err)
	}

	// Unmarshal the secret value into RDSSecretData struct
	var secretData RDSSecretData
	err = json.Unmarshal([]byte(*result.SecretString), &secretData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret value, %v", err)
	}

	return &secretData, nil
}

func main() {

	var eventLog []report.Event
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stderr).With().Str("version", version.Version).Timestamp().Logger()
	logger.Info().Msg("starting some stuff")

	//secretID := "SandboxSharedRdsInstanceMas-JyJHrRRpi8Ex"
	//secret, err := GetRDSSecret(secretID)
	//logger.Info().Msgf("host: %s", secret.Host)

	//err = CheckTCPConnectivity(*secret)
	//if err != nil {
	//	panic(err)
	//}
	//logger.Info().Msg("tcp connectivity check passed")
	eventLog = append(eventLog, report.Event{
		Time:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "Start hardware upgrade",
	})
	eventLog = append(eventLog, report.Event{
		Time:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "Succeeded",
	})
	eventLog = append(eventLog, report.Event{
		Time:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "Succeeded",
	})
	eventLog = append(eventLog, report.Event{
		Time:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "Failed",
	})
	eventLog = append(eventLog, report.Event{
		Time:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "Failed",
	})
	eventLog = append(eventLog, report.Event{
		Time:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "Succeeded",
	})
	eventLog = append(eventLog, report.Event{
		Time:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "Succeeded",
	})
	eventLog = append(eventLog, report.Event{
		Time:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "End hardware upgrade",
	})
}
