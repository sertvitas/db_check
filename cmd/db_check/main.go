package main

import (
	"os"

	"github.com/sertvitas/db_check/version"

	"github.com/rs/zerolog"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stderr).With().Str("version", version.Version).Timestamp().Logger()
	logger.Info().Msg("starting some other shit")

}
