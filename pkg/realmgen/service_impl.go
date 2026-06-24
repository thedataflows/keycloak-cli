package realmgen

import "github.com/thedataflows/keycloak-cli/pkg/realmgen/internal"

type service struct {
	inner *internal.Service
}

func newService() Service {
	return &service{inner: internal.New()}
}

func ValidateOptions(options Options) error {
	return internal.ValidateOptions(internal.Options(options))
}

func (s *service) Generate(specPath string, options Options) (Result, error) {
	result, err := s.inner.Generate(specPath, internal.Options(options))
	if err != nil {
		return Result{}, err
	}

	return Result{
		Resources:     result.Resources,
		Relationships: result.Relationships,
		Summary:       Summary(result.Summary),
	}, nil
}
