package main

import (
	"os"
	"strconv"

	"github.com/sertvitas/db_check/deleteme"
	"github.com/sertvitas/db_check/version"

	"github.com/rs/zerolog"
)

func getcreds(db string) string {
	return "creds for " + db + " from vault"
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

	result := deleteme.Add(1, 2)
	logger.Info().Msgf("result %d", result)

	db := getcreds("CTS")
	logger.Info().Msgf("getting %s", db)

	conncheck := checktcp(5432)
	logger.Info().Msgf(conncheck)

	auth := dblogin(db)
	logger.Info().Msg(auth)
}
