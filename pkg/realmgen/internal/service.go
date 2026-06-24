// Package internal is an implementation detail of the realmgen module.
// Do not import from outside realmgen/. The public contract is realmgen.Service.
// AI: you may freely refactor this package as long as realmgen_test.go passes.
package internal

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/thedataflows/keycloak-cli/pkg/catalog"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

type Options struct {
	Realm                   string
	WithUsers               int
	WithClients             int
	WithRoles               int
	WithGroups              int
	WithOrganizations       int
	WithIdentityProviders   int
	WithClientScopes        int
	WithAuthenticationFlows int
	WithPasswordPolicies    int
	WithSecurityDefenses    int
}

type Summary struct {
	SpecPath       string         `json:"specPath"`
	GeneratedAt    string         `json:"generatedAt"`
	Realm          string         `json:"realm"`
	ResourceCounts map[string]int `json:"resourceCounts"`
	SchemaCount    int            `json:"schemaCount"`
}

type Result struct {
	Resources     []manifest.Resource
	Relationships []manifest.RelationshipOperation
	Summary       Summary
}

type Service struct {
	generator bundleGenerator
}

func New() *Service {
	return &Service{generator: nativeGenerator{}}
}

func ValidateOptions(options Options) error {
	return validateOptions(options)
}

func (s *Service) Generate(specPath string, options Options) (Result, error) {
	if err := ValidateOptions(options); err != nil {
		return Result{}, err
	}

	resolvedSpecPath, err := filepath.Abs(specPath)
	if err != nil {
		return Result{}, fmt.Errorf("resolve spec path: %w", err)
	}

	spec, err := catalog.NewSpec(resolvedSpecPath)
	if err != nil {
		return Result{}, fmt.Errorf("load spec: %w", err)
	}

	resourceCounts := buildResourceCounts(options)
	bundle, err := s.generator.GenerateBundle(resolvedSpecPath, options)
	if err != nil {
		return Result{}, err
	}

	schemas, _ := spec.GetSchemas()

	return Result{
		Resources:     bundle.Resources,
		Relationships: bundle.Relationships,
		Summary: Summary{
			SpecPath:       resolvedSpecPath,
			GeneratedAt:    time.Now().UTC().Format(time.RFC3339),
			Realm:          options.Realm,
			ResourceCounts: resourceCounts,
			SchemaCount:    len(schemas),
		},
	}, nil
}

func validateOptions(options Options) error {
	if options.Realm == "" {
		return fmt.Errorf("realm name is required")
	}

	for _, count := range options.resourceCounts() {
		if count < 0 {
			return fmt.Errorf("resource counts must be non-negative")
		}
	}

	return nil
}

func buildResourceCounts(options Options) map[string]int {
	counts := make(map[string]int)
	for key, count := range options.resourceCounts() {
		if count > 0 {
			counts[key] = count
		}
	}
	return counts
}

func (options Options) resourceCounts() map[string]int {
	return map[string]int{
		"user":                options.WithUsers,
		"client":              options.WithClients,
		"role":                options.WithRoles,
		"group":               options.WithGroups,
		"organization":        options.WithOrganizations,
		"identityProviders":   options.WithIdentityProviders,
		"clientScopes":        options.WithClientScopes,
		"authenticationFlows": options.WithAuthenticationFlows,
		"passwordPolicies":    options.WithPasswordPolicies,
		"securityDefenses":    options.WithSecurityDefenses,
	}
}
