package main

import (
	"fmt"
	"os"
	"time"

	CheckTCPResponse "github.com/sertvitas/db_check/cmd/checkTCP"

	"github.com/sertvitas/db_check/cmd/rdsip"

	"github.com/sertvitas/db_check/version"

	"github.com/rs/zerolog"
)

func getcreds(db string) string {
	return "creds for " + db + " from vault"
}

//func checktcp(port int) string {
//	return "checking tcp port " + strconv.Itoa(port)
//}

func dblogin(auth string) string {
	return "logging into CTS db with " + auth
}

// func dbping()
func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stderr).With().Str("version", version.Version).Timestamp().Logger()
	logger.Info().Msg("starting some other shit")

	db := getcreds("CTS")
	logger.Info().Msgf("getting %s", db)

	ipToCheck := "212.58.249.144"
	portToCheck := 443
	timeoutDuration := 5 * time.Second
	conncheck := CheckTCPResponse.CheckTCPResponse(ipToCheck, portToCheck, timeoutDuration)
	logger.Info().Msgf("TCP response: ", conncheck)

	instanceName := "cts-cts-sandbox"
	region := "us-east-1b"
	ip, err := rdsip.GetRDSInstanceIP(instanceName, region)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	logger.Info().Msgf(ip)

	//auth := dblogin(db)
	//logger.Info().Msg(auth)
}
