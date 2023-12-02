package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/sertvitas/db_check/poll"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog/log"
	"github.com/sertvitas/db_check/report"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/sertvitas/db_check/version"

	"github.com/rs/zerolog"
)

const (
	instanceMonitorInitialDelaySeconds = 120
	instanceMonitorPollDelaySeconds    = 15
	checkDelaySeconds                  = 1
)

var instanceIsAvailable = false

var eventLog []report.Event

// CheckTCPConnectivity checks TCP connectivity to the specified host and port
func CheckTCPConnectivity(secret poll.RDSSecretData) error {
	address := fmt.Sprintf("%s:%d", secret.Host, secret.Port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", address, err)
	}
	defer conn.Close()
	return nil
}

// GetRDSSecret retrieves the secret value for the given secret ID and unmarshals it into RDSSecretData
func GetRDSSecret(secretID string) (*poll.RDSSecretData, error) {
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
	var secretData poll.RDSSecretData
	err = json.Unmarshal([]byte(*result.SecretString), &secretData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret value, %v", err)
	}

	return &secretData, nil
}

func runTests(checkDelaySeconds int, secret poll.RDSSecretData, eLog *[]report.Event) {
	for {
		if instanceIsAvailable {
			log.Info().Msg("instance is available. stopping connection tests")
			break
		}
		err := CheckTCPConnectivity(secret)
		if err != nil {
			log.Info().Msgf("TCP connectivity check failed: %v", err)
			*eLog = append(*eLog, report.Event{
				Time: time.Now(), Description: "Failed"})
		} else {
			log.Info().Msg("TCP connectivity check passed")
			*eLog = append(*eLog, report.Event{Time: time.Now(), Description: "Succeeded"})
		}
	}

}

func main() {
	// set up logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stderr).With().Str("version", version.Version).Timestamp().Logger()

	// get secret value
	secret, err := GetRDSSecret("SandboxSharedRdsInstanceMas-JyJHrRRpi8Ex")
	if err != nil {
		panic(err)
	}
	logger.Info().Msgf("got secret for: %s", secret.Host)

	//  bail out if we can't connect ot the instance before we even start the upgrade
	err = CheckTCPConnectivity(*secret)
	if err != nil {
		logger.Fatal().Err(err).Msg("TCP connectivity check failed")
	}

	// set the initial event log
	eventLog = append(eventLog, report.Event{Time: time.Now(), Description: "Start hardware upgrade"})

	// kick off the background instance status monitor.
	// it sets the instanceIsAvailable flag to true when the instance is available
	go poll.InstanceMonitor(
		&instanceIsAvailable, instanceMonitorInitialDelaySeconds, instanceMonitorPollDelaySeconds, *secret, &logger)

	// run tests until instanceIsAvailable is true
	runTests(checkDelaySeconds, *secret, &eventLog)

	// add final event to event log
	eventLog = append(eventLog, report.Event{Time: time.Now(), Description: "Hardware upgrade complete"})

	// generate report
	report := report.GetReport(eventLog)
	fmt.Println(report)

}
