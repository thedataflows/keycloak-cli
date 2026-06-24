package realmgen

import (
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

// Service is the public API for realm generation.
type Service interface {
	Generate(specPath string, options Options) (Result, error)
}

// New returns a production-ready realm generation service.
func New() Service {
	return newService()
}
