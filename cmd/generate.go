package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/thedataflows/keycloak-cli/pkg/output"
	"github.com/thedataflows/keycloak-cli/pkg/realmgen"
)

// GenerateCmd generates a complete Keycloak realm configuration with all resources
type GenerateCmd struct {
	Output                  string `short:"o" help:"Output file or directory (ends with / for directory)"`
	Format                  string `short:"f" default:"json" enum:"json,yaml,toml" help:"Output format"`
	Realm                   string `short:"r" help:"Name of the realm to generate" default:"test-realm"`
	WithUsers               int    `help:"Number of users to generate" default:"1"`
	WithClients             int    `help:"Number of clients to generate" default:"1"`
	WithRoles               int    `help:"Number of roles to generate" default:"1"`
	WithGroups              int    `help:"Number of groups to generate" default:"1"`
	WithOrganizations       int    `help:"Number of organizations to generate" default:"1"`
	WithIdentityProviders   int    `help:"Number of identity providers to generate" default:"1"`
	WithClientScopes        int    `help:"Number of client scopes to generate" default:"1"`
	WithAuthenticationFlows int    `help:"Number of authentication flows to generate" default:"1"`
	WithPasswordPolicies    int    `help:"Number of password policies to generate" default:"1"`
	WithSecurityDefenses    int    `help:"Generate security defenses configuration" default:"1"`
	Summary                 bool   `help:"Generate summary" default:"false"`
	Overwrite               bool   `help:"Overwrite existing files" default:"false"`
}

// Validate validates the command options
func (c *GenerateCmd) Validate() error {
	return realmgen.ValidateOptions(c.realmgenOptions())
}

// Run executes the generate command
func (c *GenerateCmd) Run(_ *kong.Context, cli *CLI) error {
	options := c.realmgenOptions()

	service := realmgen.New()
	result, err := service.Generate(cli.SpecPath, options)
	if err != nil {
		return err
	}

	outPath := strings.TrimSpace(c.Output)
	switch {
	case outPath == "":
		if err := output.WritePayload(os.Stdout, result.Resources, c.Format); err != nil {
			return err
		}
		if len(result.Relationships) > 0 {
			fmt.Fprintln(os.Stdout)
			if err := output.WritePayload(os.Stdout, map[string]interface{}{"relationships": result.Relationships}, c.Format); err != nil {
				return err
			}
		}
	case isDirTarget(outPath):
		if err := os.MkdirAll(outPath, 0o755); err != nil {
			return fmt.Errorf("create output directory: %w", err)
		}
		if err := c.writeFile(filepath.Join(outPath, "realm."+c.Format), result.Resources); err != nil {
			return err
		}
		log.Logger.Info().Str("pkg", PKG_CMD).Msgf("generated realm: %s", filepath.Join(outPath, "realm."+c.Format))

		if len(result.Relationships) > 0 {
			if err := c.writeFile(filepath.Join(outPath, "relationships."+c.Format), map[string]interface{}{"relationships": result.Relationships}); err != nil {
				return err
			}
			log.Logger.Info().Str("pkg", PKG_CMD).Msgf("generated relationships: %s", filepath.Join(outPath, "relationships."+c.Format))
		}
	default:
		if err := c.writeFile(outPath, result.Resources); err != nil {
			return err
		}
		log.Logger.Info().Str("pkg", PKG_CMD).Msgf("generated realm: %s", outPath)
	}

	if c.Summary {
		summaryJSON, err := json.MarshalIndent(result.Summary, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal summary: %w", err)
		}
		log.Logger.Info().Str("pkg", PKG_CMD).Msgf("generated summary: %s", summaryJSON)
	}

	log.Logger.Info().Str("pkg", PKG_CMD).Msgf("generated %d resources for realm '%s'", len(result.Resources), c.Realm)
	return nil
}

func (c *GenerateCmd) writeFile(path string, payload interface{}) error {
	if !c.Overwrite {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists (use --overwrite to replace)", path)
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()
	return output.WritePayload(file, payload, c.Format)
}

func (c *GenerateCmd) realmgenOptions() realmgen.Options {
	return realmgen.Options{
		Realm:                   c.Realm,
		WithUsers:               c.WithUsers,
		WithClients:             c.WithClients,
		WithRoles:               c.WithRoles,
		WithGroups:              c.WithGroups,
		WithOrganizations:       c.WithOrganizations,
		WithIdentityProviders:   c.WithIdentityProviders,
		WithClientScopes:        c.WithClientScopes,
		WithAuthenticationFlows: c.WithAuthenticationFlows,
		WithPasswordPolicies:    c.WithPasswordPolicies,
		WithSecurityDefenses:    c.WithSecurityDefenses,
	}
}
