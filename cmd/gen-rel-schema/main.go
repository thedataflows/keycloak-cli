package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
)

type cliConfig struct {
	SpecPath string `help:"Path to the Keycloak OpenAPI specification" type:"existingfile" short:"s" default:"keycloak-oapi/26.6.2.spec.json"`
	Output   string `help:"Path to write the generated JSON schema" short:"o" default:"generated/relationships/schema.json"`
}

func main() {
	cfg := cliConfig{}
	ctx := kong.Parse(&cfg,
		kong.Name("gen-rel-schema"),
		kong.Description("Generate JSON schema validating Keycloak relationship uploads"),
	)

	spec, err := catalog.NewSpec(cfg.SpecPath)
	if err != nil {
		ctx.FatalIfErrorf(err)
	}

	templates, collectErr := catalog.CollectRelationshipTemplates(spec)
	if collectErr != nil {
		warn("parsing spec", collectErr)
	}

	schema, buildErr := buildRelationshipJSONSchema(templates)
	if buildErr != nil {
		warn("preparing schema", buildErr)
	}

	if err := writeSchema(cfg.Output, schema); err != nil {
		ctx.FatalIfErrorf(err)
	}
}

func warn(stage string, err error) {
	fmt.Fprintf(os.Stderr, "encountered issues while %s: %v\n", stage, err)
}

func buildRelationshipJSONSchema(templates []catalog.RelationshipTemplate) (map[string]interface{}, error) {
	relationships := map[string]interface{}{
		"type": "array",
	}

	variants, variantErr := buildRelationshipVariants(templates)
	if len(variants) > 0 {
		relationships["items"] = map[string]interface{}{
			"oneOf": variants,
		}
	}

	schema := map[string]interface{}{
		"$schema":              "http://json-schema.org/draft-07/schema#",
		"title":                "Keycloak Relationship Upload Schema",
		"type":                 "object",
		"required":             []string{"relationships"},
		"additionalProperties": false,
		"properties": map[string]interface{}{
			"relationships": relationships,
		},
	}

	return schema, variantErr
}

func buildRelationshipVariants(templates []catalog.RelationshipTemplate) ([]interface{}, error) {
	variants := make([]interface{}, 0, len(templates))
	var errs []error

	for _, tmpl := range templates {
		variant, err := buildRelationshipVariant(tmpl)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s %s: %w", tmpl.Method, tmpl.Template, err))
			continue
		}
		variants = append(variants, variant)
	}

	return variants, errors.Join(errs...)
}

func buildRelationshipVariant(tmpl catalog.RelationshipTemplate) (map[string]interface{}, error) {
	pattern, err := catalog.RelationshipTemplatePattern(tmpl.Template)
	if err != nil {
		return nil, err
	}

	properties := map[string]interface{}{
		"path": map[string]interface{}{
			"type":    "string",
			"pattern": pattern,
		},
	}

	required := []string{"path"}

	if tmpl.Method == http.MethodDelete {
		properties["delete"] = map[string]interface{}{
			"type":  "boolean",
			"const": true,
		}
		required = append(required, "delete")
	} else {
		properties["delete"] = map[string]interface{}{
			"type": "boolean",
			"enum": []interface{}{false},
		}
	}

	if tmpl.Body != nil {
		properties["data"] = tmpl.Body
	}

	if tmpl.RequiresBody {
		if _, ok := properties["data"]; !ok {
			properties["data"] = map[string]interface{}{}
		}
		required = append(required, "data")
	}

	variant := map[string]interface{}{
		"type":                 "object",
		"required":             required,
		"additionalProperties": false,
		"properties":           properties,
	}

	if tmpl.Summary != "" {
		variant["description"] = tmpl.Summary
	}

	return variant, nil
}

func writeSchema(dest string, schema map[string]interface{}) error {
	if dest == "" {
		return fmt.Errorf("output path is empty")
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return err
	}

	return nil
}
