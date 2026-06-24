package admin

import (
	"context"
	"time"

	"github.com/thedataflows/keycloak-cli/pkg/catalog"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

type Config struct {
	BaseURL  string
	SpecPath string
	Timeout  time.Duration
}

// Service is the public admin API used by command handlers.
type Service interface {
	Spec() *catalog.Spec
	Fetch(ctx context.Context, query FetchQuery) (FetchReport, error)
	Apply(ctx context.Context, resources []manifest.Resource, relationships []manifest.RelationshipOperation, options ApplyOptions) (ApplyReport, error)
}

func New(config Config) (Service, error) {
	return newService(config)
}
