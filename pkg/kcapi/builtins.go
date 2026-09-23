package kcapi

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

// BuiltInResource describes a resource type and the names that should be treated
// as built-in and ignored during round-trip comparison.
type BuiltInResource struct {
	Type  string   `yaml:"type"`
	Names []string `yaml:"names"`
}

// BuiltInResourcesFile is the top-level shape of a built-in resources YAML file.
type BuiltInResourcesFile struct {
	Resources []BuiltInResource `yaml:"resources"`
}

var defaultBuiltInResourceNames = map[string][]string{
	"client": {"account", "account-console", "admin-cli", "broker", "realm-management", "security-admin-console"},
	"role":   {"offline_access", "uma_authorization", "default-roles-*"},
	"clientscope": {
		"acr", "address", "basic", "email", "microprofile-jwt", "offline_access",
		"organization", "phone", "profile", "role_list", "roles", "saml_organization",
		"service_account", "web-origins",
	},
}

// LoadBuiltInResources reads built-in resource overrides from a YAML file.
func LoadBuiltInResources(path string) ([]BuiltInResource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read built-in resources: %w", err)
	}

	var file BuiltInResourcesFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse built-in resources: %w", err)
	}
	return file.Resources, nil
}

// ApplyBuiltInResources merges override entries into the default built-in name map.
func ApplyBuiltInResources(overrides []BuiltInResource) error {
	for _, override := range overrides {
		if override.Type == "" {
			return fmt.Errorf("built-in resource override missing required field 'type'")
		}
		defaultBuiltInResourceNames[override.Type] = patchFieldList(
			defaultBuiltInResourceNames[override.Type],
			override.Names,
			nil,
		)
	}
	return nil
}

// InstallDefaultBuiltInResources loads built-in resource overrides from the
// directory containing the spec and wires the matcher into the package-level
// hook below. It is safe to call when no override file exists.
func InstallDefaultBuiltInResources(specPath string) error {
	specDir := filepath.Dir(specPath)
	overridePath := filepath.Join(specDir, "built-in-resources.yaml")

	IsBuiltInResource = builtInResourceMatcher

	overrides, err := LoadBuiltInResources(overridePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return ApplyBuiltInResources(overrides)
}

// IsBuiltInResource is the hook var moved verbatim from pkg/manifest/manifest.go
// (manifest's round-trip reads now consult kcapi.IsBuiltInResource directly).
// It allows the catalog package to inject knowledge of built-in
// resources that should be excluded from round-trip comparison.
var IsBuiltInResource = func(Resource) bool { return false }

func builtInResourceMatcher(resource Resource) bool {
	names, ok := defaultBuiltInResourceNames[resource.Type]
	if !ok {
		return false
	}
	name := builtInResourceName(resource)
	for _, candidate := range names {
		if matchSimplePattern(name, candidate) {
			return true
		}
	}
	return false
}

func builtInResourceName(resource Resource) string {
	if name := stringField(resource.Data, "clientId"); name != "" {
		return name
	}
	if name := stringField(resource.Data, "name"); name != "" {
		return name
	}
	if name := stringField(resource.Data, "alias"); name != "" {
		return name
	}
	if name := stringField(resource.Data, "username"); name != "" {
		return name
	}
	return stringField(resource.Data, "realm")
}

func matchSimplePattern(value, pattern string) bool {
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(value, pattern[:len(pattern)-1])
	}
	if strings.HasPrefix(pattern, "*") {
		return strings.HasSuffix(value, pattern[1:])
	}
	return value == pattern
}
