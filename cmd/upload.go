package cmd

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
	"github.com/thedataflows/keycloak-cli/pkg/output"
)

// UploadCmd uploads objects to Keycloak
type UploadCmd struct {
	InputFiles      []string `arg:"" help:"Input JSON/YAML files to upload"`
	DryRun          bool     `help:"Show what would be uploaded without actually doing it" default:"false"`
	ContinueOnError bool     `help:"Continue uploading other resources if one fails" default:"false"`
	Format          string   `short:"f" default:"table" enum:"table,json,yaml" help:"Output format"`
	Output          string   `short:"o" help:"Output file or directory for upload results (ends with / for directory)"`
	Force           bool     `help:"Force overwrite of output file" default:"false"`
	Delete          bool     `help:"Delete resources regardless if they are marked for deletion or not" default:"false"`
}

func (c *UploadCmd) Run(ctx *kong.Context, cli *CLI) error {
	if len(c.InputFiles) == 0 {
		return fmt.Errorf("at least one input file or directory required")
	}

	specClient, err := cli.adminClient()
	if err != nil {
		return err
	}

	commandCtx, cancel := cli.CreateContextWithTimeout()
	defer cancel()

	loaded, err := manifest.LoadPaths(c.InputFiles)
	if err != nil {
		return err
	}
	for _, skipped := range loaded.Skipped {
		log.Logger.Warn().Str("pkg", PKG_CMD).Msgf("skipping %s: %s", skipped.Path, skipped.Reason)
	}

	report, err := specClient.Apply(commandCtx, loaded.Resources, loaded.Relationships, admin.ApplyOptions{
		DryRun:          c.DryRun,
		Delete:          c.Delete,
		ContinueOnError: c.ContinueOnError,
		Reconcile:       true,
	})
	if outputErr := c.outputResults(report.Results); outputErr != nil {
		return outputErr
	}
	if err != nil {
		return err
	}

	log.Logger.Info().Str("pkg", PKG_CMD).Msgf("uploaded %d resources (%d failed)", len(report.Results), report.Failed)
	return nil
}

func (c *UploadCmd) outputResults(results []admin.ApplyResult) error {
	dest, shouldClose, err := output.Destination(c.Output, c.Force)
	if err != nil {
		return err
	}
	if shouldClose {
		defer dest.Close()
	}

	return output.WriteApplyResults(dest, results, c.Format)
}
