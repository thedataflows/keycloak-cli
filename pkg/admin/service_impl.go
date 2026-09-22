package admin

import (
	"time"

	"github.com/thedataflows/keycloak-cli/pkg/auth"
	"github.com/thedataflows/keycloak-cli/pkg/kcapi"
)

type service struct {
	specClient *kcapi.RuntimeClient
	timeout    time.Duration

	identities    map[string]kcapi.ResourceIdentity
	identitiesErr error
}

func newService(config Config) (Service, error) {
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
