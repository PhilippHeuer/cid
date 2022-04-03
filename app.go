package main

import (
	"github.com/cidverse/cid/pkg/cmd"
	"github.com/rs/zerolog/log"
)

// Version will be set at build time
var Version string

// CommitHash will be set at build time
var CommitHash string

// BuildAt will be set at build time
var BuildAt string

// Init Hook
func init() {
	// Set Version Information
	cmd.Version = Version
	cmd.CommitHash = CommitHash
	cmd.BuildAt = BuildAt
}

// CLI Main Entrypoint
func main() {
	cmdErr := cmd.Execute()
	if cmdErr != nil {
		log.Fatal().Err(cmdErr).Msg("cli error")
	}
}
