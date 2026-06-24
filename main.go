package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/thedataflows/keycloak-cli/cmd"
)

var version = "dev"

func main() {
	err := cmd.Run(version, os.Args[1:])
	if err != nil {
		log.Logger.Error().Str("pkg", "main").Err(err).Msg("Command failed")
		os.Exit(1)
	}
}
