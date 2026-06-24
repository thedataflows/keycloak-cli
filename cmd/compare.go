package cmd

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
	"github.com/thedataflows/keycloak-cli/pkg/output"
)

type CompareCmd struct {
	InputFiles []string `arg:"" help:"Input JSON/YAML files to compare against Keycloak state"`
	Realm      string   `short:"r" help:"Realm to compare. If not provided, infer it from the manifests."`
	Format     string   `short:"f" default:"table" enum:"table,json,yaml" help:"Output format"`
	Output     string   `short:"o" help:"Output file or directory for comparison report (ends with / for directory)"`
	Force      bool     `help:"Force overwrite of output file" default:"false"`
}

func (c *CompareCmd) Run(ctx *kong.Context, cli *CLI) error {
	if len(c.InputFiles) == 0 {
		return fmt.Errorf("at least one input file or directory required")
	}

	svc, err := cli.adminClient()
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
	if len(loaded.Resources) == 0 && len(loaded.Relationships) == 0 {
		return fmt.Errorf("no manifests loaded")
	}
	if err := catalog.ValidateRelationshipOperations(svc.Spec(), loaded.Relationships); err != nil {
		return err
	}

	realm, err := c.resolveRealm(loaded)
	if err != nil {
		return err
	}

	fetched, err := svc.Fetch(commandCtx, admin.FetchQuery{
		Realm:                realm,
		Resources:            compareFetchResources(loaded),
		IncludeRelationships: len(loaded.Relationships) > 0,
	})
	if err != nil {
		return err
	}

	report := manifest.CompareRoundTrip(loaded.Resources, loaded.Relationships, fetched.Resources, fetched.Relationships)
	if err := c.outputReport(report); err != nil {
		return err
	}
	if !report.Match {
		return fmt.Errorf("comparison failed")
	}
	return nil
}

func (c *CompareCmd) outputReport(report manifest.ComparisonReport) error {
	dest, shouldClose, err := output.Destination(c.Output, c.Force)
	if err != nil {
		return err
	}
	if shouldClose {
		defer dest.Close()
	}
	return output.WriteComparisonReport(dest, report, c.Format)
}

func (c *CompareCmd) resolveRealm(loaded manifest.LoadResult) (string, error) {
	if strings.TrimSpace(c.Realm) != "" {
		return strings.TrimSpace(c.Realm), nil
	}

	realms := collectUniqueRealms(loaded)
	if len(realms) == 1 {
		return realms[0], nil
	}
	if len(realms) == 0 {
		return "", fmt.Errorf("realm is required")
	}
	return "", fmt.Errorf("multiple realms found in input; use --realm")
}

func collectUniqueRealms(loaded manifest.LoadResult) []string {
	seen := make(map[string]struct{})
	realms := make([]string, 0)
	add := func(realm string) {
		realm = strings.TrimSpace(realm)
		if realm == "" {
			return
		}
		if _, exists := seen[realm]; exists {
			return
		}
		seen[realm] = struct{}{}
		realms = append(realms, realm)
	}

	for _, resource := range loaded.Resources {
		realm := resource.Realm
		if realm == "" {
			realm = resource.Name()
		}
		add(realm)
	}
	for _, relationship := range loaded.Relationships {
		parts := strings.Split(strings.TrimPrefix(strings.TrimSpace(relationship.Path), "/"), "/")
		if len(parts) > 0 {
			add(parts[0])
		}
	}
	return realms
}

func compareFetchResources(loaded manifest.LoadResult) string {
	resourceSet := make(map[string]struct{})
	for _, resource := range loaded.Resources {
		if resource.Type != "" {
			resourceSet[resource.Type] = struct{}{}
		}
	}
	for _, rel := range loaded.Relationships {
		for _, resourceType := range manifest.RelationshipParamTypes(rel.Kind) {
			resourceSet[resourceType] = struct{}{}
		}
	}
	if len(resourceSet) == 0 {
		resourceSet["realm"] = struct{}{}
	}
	resources := slices.Sorted(maps.Keys(resourceSet))
	return strings.Join(resources, ",")
}
