package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/natemarks/secret-hoard/types"
	"github.com/sertvitas/db_check/version"

	"github.com/rs/zerolog"
)

const (
	// MasterSecretEnvVar is the name of the environment variable that contains master secret value
	MasterSecretEnvVar = "MASTER_SECRET"
)

//func getcreds(db string) string {
//	return "creds for " + db + " from vault"
//}

func getSecretFromEnvVar(envVarKey string, log *zerolog.Logger) (secret types.RDSSecretData, err error) {
	secretStr, err := getEnvVar(envVarKey)
	if err != nil {
		log.Error().Err(err).Msgf("error getting env var %s", envVarKey)
		return secret, err
	}

	err = json.Unmarshal([]byte(secretStr), &secret)
	if err != nil {
		log.Error().Err(err).Msgf("error unmarshalling secret from env var %s", envVarKey)
		return secret, err
	}
	log.Info().Msgf("secret from env var %s: %v", envVarKey, secret)
	return secret, err
}

func getEnvVar(key string) (value string, err error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return value, fmt.Errorf("Environment variable '%s' not set", key)
	}
	return value, nil
}

func checktcp(port int) string {
	return "checking tcp port " + strconv.Itoa(port)
}

func dblogin(auth string) string {
	return "logging into CTS db with " + auth
}

// func dbping()
func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stderr).With().Str("version", version.Version).Timestamp().Logger()
	logger.Info().Msg("starting some other shit")

	db, _ := getSecretFromEnvVar("CTS", &logger)
	logger.Info().Msgf("getting %s", db)

	conncheck := checktcp(5432)
	logger.Info().Msgf(conncheck)

	//auth := dblogin(db)
	//logger.Info().Msg(auth)
}
