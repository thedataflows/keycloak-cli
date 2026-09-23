package manifest

import (
	"context"
	"time"

	"github.com/thedataflows/keycloak-cli/pkg/auth"
	"github.com/thedataflows/keycloak-cli/pkg/kcapi"
)

type Config struct {
	BaseURL  string
	SpecPath string
	Timeout  time.Duration

	// Auth optionally injects a custom auth.Service. When nil, NewService
	// falls back to auth.New() (the password-grant default), preserving
	// backward compatibility. The syncengine migration injects a
	// client_credentials provider here.
	Auth auth.Service
}

// Service is the public manifest API used by command handlers.
type Service interface {
	Spec() *kcapi.Spec
	Fetch(ctx context.Context, query FetchQuery) (FetchReport, error)
	// FetchChildren returns one child collection of one parent resource (e.g. the
	// roles of a client) with exactly one HTTP GET — no depth fan-out and no
	// realm-wide reference-resolution sweep. See fetch.go for the contract.
	FetchChildren(ctx context.Context, parent Resource, childType string, query ChildFetchQuery) (FetchReport, error)
	Apply(ctx context.Context, resources []Resource, relationships []RelationshipOperation, options ApplyOptions) (ApplyReport, error)
}

type service struct {
	specClient *kcapi.RuntimeClient
	timeout    time.Duration

	identities    map[string]kcapi.ResourceIdentity
	identitiesErr error
}

// NewService builds a manifest Service on top of a kcapi client: the runtime
// client is constructed first, then wrapped in the thin service below.
func NewService(config Config) (Service, error) {
	var authSvc auth.Service = auth.New()
	if config.Auth != nil {
		authSvc = config.Auth
	}
	impl, err := kcapi.NewRuntimeClient(kcapi.RuntimeConfig{
		BaseURL:  config.BaseURL,
		SpecPath: config.SpecPath,
		Timeout:  config.Timeout,
	}, authSvc)
	if err != nil {
		return nil, err
	}

	return &service{specClient: impl, timeout: config.Timeout}, nil
}

func (s *service) Spec() *kcapi.Spec {
	return s.specClient.Spec()
}

func (s *service) resourceIdentity(resourceType string) (kcapi.ResourceIdentity, bool) {
	if s.identities == nil && s.identitiesErr == nil {
		s.identities, s.identitiesErr = s.Spec().ResourceIdentities()
	}
	identity, ok := s.identities[resourceType]
	return identity, ok
}
