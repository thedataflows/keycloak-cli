package catalog

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/thedataflows/keycloak-cli/pkg/catalog/internal"
)

type Spec struct {
	internal *internal.Spec
}

func NewSpec(path string) (*Spec, error) {
	loaded, err := internal.NewSpec(path)
	if err != nil {
		return nil, err
	}
	return &Spec{internal: loaded}, nil
}

func NewSpecFromBytes(data []byte) (*Spec, error) {
	loaded, err := internal.NewSpecFromBytes(data)
	if err != nil {
		return nil, err
	}
	return &Spec{internal: loaded}, nil
}

func (s *Spec) GetSchemas() (map[string]*base.SchemaProxy, error) {
	if s == nil {
		return map[string]*base.SchemaProxy{}, nil
	}
	return s.internal.GetSchemas()
}

func (s *Spec) ForEachOperation(visitor func(path, method string, operation *v3.Operation, item *v3.PathItem)) {
	if s == nil {
		return
	}
	s.internal.ForEachOperation(visitor)
}

func (s *Spec) Operation(path, method string) (*v3.Operation, *v3.PathItem, error) {
	if s == nil {
		return nil, nil, internal.ErrSpecNotInitialized
	}
	return s.internal.Operation(path, method)
}

func (s *Spec) PathParameters(path, method string) ([]string, error) {
	if s == nil {
		return nil, internal.ErrSpecNotInitialized
	}
	return s.internal.PathParameters(path, method)
}

// WrapSpec exposes an already-loaded internal spec for tests that embed a fixture.
func WrapSpec(loaded *internal.Spec) *Spec {
	return &Spec{internal: loaded}
}
