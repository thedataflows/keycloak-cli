// Package internal is an implementation detail of the realmgen module.
// Do not import from outside realmgen/. The public contract is realmgen.Service.
// AI: you may freely refactor this package as long as realmgen_test.go passes.
package internal

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

type generatedBundle struct {
	Resources     []manifest.Resource
	Relationships []manifest.RelationshipOperation
}

type bundleGenerator interface {
	GenerateBundle(specPath string, options Options) (generatedBundle, error)
}

type nativeGenerator struct{}

func (nativeGenerator) GenerateBundle(specPath string, options Options) (generatedBundle, error) {
	bundle := generatedBundle{
		Resources:     make([]manifest.Resource, 0),
		Relationships: make([]manifest.RelationshipOperation, 0),
	}

	spec, err := catalog.NewSpec(specPath)
	if err != nil {
		return generatedBundle{}, fmt.Errorf("load spec: %w", err)
	}

	realm := manifest.Resource{
		Type:  "realm",
		Realm: options.Realm,
		Data: map[string]interface{}{
			"realm":        options.Realm,
			"enabled":      true,
			"displayName":  options.Realm,
			"sslRequired":  "external",
			"loginTheme":   "keycloak",
			"accountTheme": "keycloak",
			"adminTheme":   "keycloak",
			"emailTheme":   "keycloak",
		},
	}
	if options.WithOrganizations > 0 {
		realm.Data["organizationsEnabled"] = true
	}
	if options.WithPasswordPolicies > 0 {
		realm.Data["passwordPolicy"] = "length(12) and digits(1) and upperCase(1) and lowerCase(1)"
	}
	if options.WithSecurityDefenses > 0 {
		realm.Data["bruteForceProtected"] = true
	}
	bundle.Resources = append(bundle.Resources, realm)

	users := make([]manifest.Resource, 0, options.WithUsers)
	for i := range options.WithUsers {
		id := generateUUID()
		users = append(users, manifest.Resource{
			Type:  "user",
			Realm: options.Realm,
			Data: map[string]interface{}{
				"id":            id,
				"username":      fmt.Sprintf("user-%d", i+1),
				"email":         fmt.Sprintf("user-%d@example.com", i+1),
				"firstName":     fmt.Sprintf("User%d", i+1),
				"lastName":      "Generated",
				"enabled":       true,
				"emailVerified": true,
				"credentials": []map[string]interface{}{{
					"type":      "password",
					"value":     generateSecret(),
					"temporary": false,
				}},
			},
		})
	}
	bundle.Resources = append(bundle.Resources, users...)

	clients := generateNamedResources(options.Realm, "client", options.WithClients, func(i int) map[string]interface{} {
		return map[string]interface{}{
			"id":                        generateUUID(),
			"clientId":                  fmt.Sprintf("client-%d", i),
			"name":                      fmt.Sprintf("Client %d", i),
			"enabled":                   true,
			"protocol":                  "openid-connect",
			"publicClient":              false,
			"directAccessGrantsEnabled": true,
		}
	})
	bundle.Resources = append(bundle.Resources, clients...)

	roles := generateNamedResources(options.Realm, "role", options.WithRoles, func(i int) map[string]interface{} {
		return map[string]interface{}{
			"id":          generateUUID(),
			"name":        fmt.Sprintf("role-%d", i),
			"description": fmt.Sprintf("Generated role %d", i),
			"composite":   false,
			"clientRole":  false,
		}
	})
	bundle.Resources = append(bundle.Resources, roles...)

	groups := generateNamedResources(options.Realm, "group", options.WithGroups, func(i int) map[string]interface{} {
		name := fmt.Sprintf("group-%d", i)
		return map[string]interface{}{
			"id":          generateUUID(),
			"name":        name,
			"path":        "/" + name,
			"description": fmt.Sprintf("Generated group %d", i),
		}
	})
	bundle.Resources = append(bundle.Resources, groups...)

	bundle.Resources = append(bundle.Resources, generateNamedResources(options.Realm, "organization", options.WithOrganizations, func(index int) map[string]interface{} {
		return map[string]interface{}{
			"id":          generateUUID(),
			"name":        fmt.Sprintf("organization-%d", index),
			"alias":       fmt.Sprintf("org-%d", index),
			"enabled":     true,
			"domains":     []map[string]interface{}{{"name": fmt.Sprintf("org-%d.example.com", index)}},
			"description": fmt.Sprintf("Generated organization %d", index),
		}
	})...)

	bundle.Resources = append(bundle.Resources, generateNamedResources(options.Realm, "identityprovider", options.WithIdentityProviders, func(index int) map[string]interface{} {
		return map[string]interface{}{
			"id":          generateUUID(),
			"alias":       fmt.Sprintf("idp-%d", index),
			"displayName": fmt.Sprintf("Identity Provider %d", index),
			"providerId":  "oidc",
			"enabled":     true,
			"config": map[string]interface{}{
				"clientId":     fmt.Sprintf("idp-client-%d", index),
				"clientSecret": generateSecret(),
			},
		}
	})...)

	bundle.Resources = append(bundle.Resources, generateNamedResources(options.Realm, "clientscope", options.WithClientScopes, func(index int) map[string]interface{} {
		return map[string]interface{}{
			"id":          generateUUID(),
			"name":        fmt.Sprintf("scope-%d", index),
			"protocol":    "openid-connect",
			"description": fmt.Sprintf("Generated client scope %d", index),
		}
	})...)

	bundle.Resources = append(bundle.Resources, generateNamedResources(options.Realm, "authenticationflow", options.WithAuthenticationFlows, func(index int) map[string]interface{} {
		return map[string]interface{}{
			"id":          generateUUID(),
			"alias":       fmt.Sprintf("flow-%d", index),
			"description": fmt.Sprintf("Generated authentication flow %d", index),
			"providerId":  "basic-flow",
			"topLevel":    true,
			"builtIn":     false,
		}
	})...)

	bundle.Relationships = append(bundle.Relationships, generateUserGroupRelationships(options.Realm, users, groups)...)
	bundle.Relationships = append(bundle.Relationships, generateUserRoleRelationships(options.Realm, users, roles)...)
	bundle.Relationships = append(bundle.Relationships, generateGroupRoleRelationships(options.Realm, groups, roles)...)
	bundle.Relationships = append(bundle.Relationships, generateOrganizationRelationships(options.Realm, users, bundle.Resources)...)

	applySpecDefaults(spec, bundle.Resources)
	if err := catalog.ValidateResources(spec, bundle.Resources, false); err != nil {
		return generatedBundle{}, fmt.Errorf("validate generated resources: %w", err)
	}
	if err := catalog.ValidateRelationshipOperations(spec, bundle.Relationships); err != nil {
		return generatedBundle{}, fmt.Errorf("validate generated relationships: %w", err)
	}

	return bundle, nil
}

func applySpecDefaults(spec *catalog.Spec, resources []manifest.Resource) {
	if spec == nil {
		return
	}
	contracts, err := spec.ResourceContracts()
	if err != nil {
		return
	}
	for idx := range resources {
		contract, ok := contracts[resources[idx].Type]
		if !ok || contract.Schema == nil {
			continue
		}
		applySchemaDefaults(contract.Schema, resources[idx].Data, resources[idx].Type)
	}
}

func applySchemaDefaults(proxy *base.SchemaProxy, data map[string]interface{}, seed string) {
	if proxy == nil || data == nil {
		return
	}
	schema, err := proxy.BuildSchema()
	if err != nil || schema == nil {
		return
	}
	for _, required := range schema.Required {
		if _, exists := data[required]; exists {
			continue
		}
		data[required] = schemaPlaceholder(schemaProperty(schema, required), required, seed)
	}
	if schema.Properties == nil {
		return
	}
	for pair := schema.Properties.First(); pair != nil; pair = pair.Next() {
		child, exists := data[pair.Key()]
		if !exists {
			continue
		}
		childMap, ok := child.(map[string]interface{})
		if !ok {
			continue
		}
		applySchemaDefaults(pair.Value(), childMap, pair.Key())
	}
}

func schemaProperty(schema *base.Schema, name string) *base.SchemaProxy {
	if schema == nil || schema.Properties == nil {
		return nil
	}
	return schema.Properties.GetOrZero(name)
}

func schemaPlaceholder(proxy *base.SchemaProxy, field, seed string) interface{} {
	if proxy == nil {
		return defaultSeedValue(field, seed)
	}
	schema, err := proxy.BuildSchema()
	if err != nil || schema == nil {
		return defaultSeedValue(field, seed)
	}
	if len(schema.Enum) > 0 && schema.Enum[0] != nil {
		var value interface{}
		if schema.Enum[0].Decode(&value) == nil {
			return value
		}
	}
	if len(schema.Type) == 0 {
		return defaultSeedValue(field, seed)
	}
	switch schema.Type[0] {
	case "boolean":
		return true
	case "integer":
		return 1
	case "number":
		return 1
	case "array":
		if schema.Items != nil && schema.Items.IsA() && schema.Items.A != nil {
			return []interface{}{schemaPlaceholder(schema.Items.A, field, seed)}
		}
		return []interface{}{}
	case "object":
		if strings.EqualFold(field, "config") {
			return map[string]interface{}{}
		}
		mapped := map[string]interface{}{}
		for _, required := range schema.Required {
			mapped[required] = schemaPlaceholder(schemaProperty(schema, required), required, field)
		}
		return mapped
	default:
		return defaultSeedValue(field, seed)
	}
}

func defaultSeedValue(field, seed string) interface{} {
	switch strings.ToLower(field) {
	case "name", "username", "alias", "realm", "clientid", "providerid", "provider", "type":
		if seed == "" {
			return field
		}
		return seed + "-" + field
	case "enabled", "publicclient", "directaccessgrantsenabled", "emailverified", "temporary", "toplevel", "builtin":
		return true
	case "description", "displayname":
		if seed == "" {
			return field
		}
		return seed + " " + field
	case "path":
		if seed == "" {
			return "/generated"
		}
		return "/" + seed
	default:
		return field
	}
}

func generateNamedResources(realm, resourceType string, count int, build func(index int) map[string]interface{}) []manifest.Resource {
	resources := make([]manifest.Resource, 0, count)
	for i := range count {
		resources = append(resources, manifest.Resource{Type: resourceType, Realm: realm, Data: build(i + 1)})
	}
	return resources
}

func generateUserGroupRelationships(realm string, users, groups []manifest.Resource) []manifest.RelationshipOperation {
	return generatePairRelationships(realm, users, groups,
		"/admin/realms/{realm}/users/{user-id}/groups/{groupId}", http.MethodPut,
		"user-id", "groupId",
		func(_, _ manifest.Resource) interface{} { return nil })
}

func generateUserRoleRelationships(realm string, users, roles []manifest.Resource) []manifest.RelationshipOperation {
	return generatePairRelationships(realm, users, roles,
		"/admin/realms/{realm}/users/{user-id}/role-mappings/realm", http.MethodPost,
		"user-id", "",
		rolePayload)
}

func generateGroupRoleRelationships(realm string, groups, roles []manifest.Resource) []manifest.RelationshipOperation {
	return generatePairRelationships(realm, groups, roles,
		"/admin/realms/{realm}/groups/{group-id}/role-mappings/realm", http.MethodPost,
		"group-id", "",
		rolePayload)
}

func generateOrganizationRelationships(realm string, users []manifest.Resource, resources []manifest.Resource) []manifest.RelationshipOperation {
	organizations := make([]manifest.Resource, 0)
	for _, resource := range resources {
		if resource.Type == "organization" {
			organizations = append(organizations, resource)
		}
	}
	if len(users) == 0 || len(organizations) == 0 {
		return nil
	}

	relationships := make([]manifest.RelationshipOperation, 0, len(organizations))
	for i, organization := range organizations {
		user := users[i%len(users)]
		orgID, _ := organization.Data["id"].(string)
		userID, _ := user.Data["id"].(string)
		if orgID == "" || userID == "" {
			continue
		}
		relationship, err := manifest.NewRelationshipOperation(
			"/admin/realms/{realm}/organizations/{org-id}/members",
			http.MethodPost,
			map[string]string{"realm": realm, "org-id": orgID},
			userID,
		)
		if err == nil {
			relationships = append(relationships, relationship)
		}
	}
	return relationships
}

func generatePairRelationships(
	realm string,
	sources, targets []manifest.Resource,
	pathTemplate, method string,
	sourceParam, targetParam string,
	payload func(source, target manifest.Resource) interface{},
) []manifest.RelationshipOperation {
	if len(sources) == 0 || len(targets) == 0 {
		return nil
	}
	relationships := make([]manifest.RelationshipOperation, 0, len(sources))
	for i, source := range sources {
		target := targets[i%len(targets)]
		params := map[string]string{"realm": realm}
		if sourceID, ok := source.Data["id"].(string); ok && sourceID != "" {
			params[sourceParam] = sourceID
		} else {
			continue
		}
		if targetParam != "" {
			if targetID, ok := target.Data["id"].(string); ok && targetID != "" {
				params[targetParam] = targetID
			} else {
				continue
			}
		}
		relationship, err := manifest.NewRelationshipOperation(
			pathTemplate,
			method,
			params,
			payload(source, target),
		)
		if err == nil {
			relationships = append(relationships, relationship)
		}
	}
	return relationships
}

func rolePayload(_, role manifest.Resource) interface{} {
	roleID, _ := role.Data["id"].(string)
	roleName, _ := role.Data["name"].(string)
	if roleID == "" || roleName == "" {
		return nil
	}
	return []map[string]interface{}{{"id": roleID, "name": roleName}}
}

func generateSecret() string {
	const lower = "abcdefghijklmnopqrstuvwxyz"
	const upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const digits = "0123456789"
	const all = lower + upper + digits

	buf := make([]byte, 32)
	rnd := make([]byte, len(buf))
	if _, err := rand.Read(rnd); err != nil {
		now := time.Now().UnixNano()
		for i := range buf {
			buf[i] = all[int(now+int64(i))%len(all)]
		}
	} else {
		for i := range buf {
			buf[i] = all[int(rnd[i])%len(all)]
		}
	}
	// Ensure policy compliance: at least one of each required class.
	buf[0] = lower[int(rnd[0])%len(lower)]
	buf[1] = upper[int(rnd[1])%len(upper)]
	buf[2] = digits[int(rnd[2])%len(digits)]
	return string(buf)
}

func generateUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return fmt.Sprintf("generated-%d", len(uuid))
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}
