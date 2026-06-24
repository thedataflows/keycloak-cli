package admin_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
	"github.com/thedataflows/keycloak-cli/pkg/auth"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestIntegrationApplyFetchAndRelationships(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test skipped in short mode")
	}
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker is not available")
	}

	startIntegrationStack(t)
	loadEnv(t)
	bootstrapAdminToken(t)

	baseURL := os.Getenv("KEYCLOAK_BASE_URL")
	require.NotEmpty(t, baseURL)

	svc, err := admin.New(admin.Config{
		BaseURL:  baseURL,
		SpecPath: filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"),
		Timeout:  30 * time.Second,
	})
	require.NoError(t, err)

	realmName := fmt.Sprintf("it-%d", time.Now().UnixNano())
	resources := []manifest.Resource{
		{Type: "realm", Realm: realmName, Data: map[string]interface{}{"realm": realmName, "enabled": true}},
		{Type: "clientscope", Realm: realmName, Data: map[string]interface{}{"name": "profile", "protocol": "openid-connect"}},
		{Type: "group", Realm: realmName, Data: map[string]interface{}{"name": "devs"}},
		{Type: "role", Realm: realmName, Data: map[string]interface{}{"name": "developer", "clientRole": false, "composite": false}},
		{Type: "client", Realm: realmName, Data: map[string]interface{}{"clientId": "app", "protocol": "openid-connect", "defaultClientScopes": []interface{}{"profile"}}},
		{Type: "user", Realm: realmName, Data: map[string]interface{}{"username": "alice", "enabled": true, "emailVerified": true, "credentials": []interface{}{map[string]interface{}{"type": "password", "value": "Password123!", "temporary": false}}}},
	}
	// First apply should create everything cleanly (including relationship ID rewriting).
	report, err := svc.Apply(context.Background(), resources, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	assert.Zero(t, report.Failed)
	// Re-apply the same manifest should report 0 failures (idempotent).
	reapplyReport, err := svc.Apply(context.Background(), resources, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	assert.Zero(t, reapplyReport.Failed)
	fetched, err := svc.Fetch(context.Background(), admin.FetchQuery{Realm: realmName, Resources: "realm,user,client,group,role,clientscope", IncludeRelationships: true})
	require.NoError(t, err)
	assert.NotEmpty(t, fetched.Resources)
	// CompareRoundTrip should match after normalization (credentials stripped, volatile fields removed).
	comparison := manifest.CompareRoundTrip(resources, nil, fetched.Resources, fetched.Relationships)
	assert.True(t, comparison.Match, "first apply + fetch should round-trip cleanly")
	assert.Empty(t, comparison.MismatchedResources)
	assert.Empty(t, comparison.MissingResources)
	assert.Empty(t, comparison.UnexpectedResources)
	lookup := make(map[string]manifest.Resource)
	for _, resource := range fetched.Resources {
		lookup[resource.Type+":"+resource.DisplayName()] = resource
	}
	user := lookup["user:alice"]
	group := lookup["group:devs"]
	role := lookup["role:developer"]
	require.NotEmpty(t, user.Identifier())
	require.NotEmpty(t, group.Identifier())
	require.NotEmpty(t, role.Identifier())
	relationships := []manifest.RelationshipOperation{
		{Path: fmt.Sprintf("%s/users/%s/groups/%s", realmName, user.Identifier(), group.Identifier())},
	}
	roleBody := []map[string]interface{}{{"id": role.Identifier(), "name": role.Name()}}
	userRole, err := manifest.NewRelationshipOperation("/admin/realms/{realm}/users/{user-id}/role-mappings/realm", "POST", map[string]string{"realm": realmName, "user-id": user.Identifier()}, roleBody)
	require.NoError(t, err)
	relationships = append(relationships, userRole)
	applyRelationships, err := svc.Apply(context.Background(), nil, relationships, admin.ApplyOptions{})
	require.NoError(t, err)
	assert.Zero(t, applyRelationships.Failed)
	fetchedWithRelationships, err := svc.Fetch(context.Background(), admin.FetchQuery{Realm: realmName, Resources: "realm,user,client,group,role,clientscope", IncludeRelationships: true})
	require.NoError(t, err)
	assert.NotEmpty(t, fetchedWithRelationships.Relationships)
	assertRoundTripResourceSubset(t, resources, fetchedWithRelationships.Resources)
	assertRoundTripRelationshipsEqual(t, fetchedWithRelationships.Resources, relationships, fetchedWithRelationships.Relationships)
	assert.Contains(t, relationshipKinds(fetchedWithRelationships.Relationships), "user-group-membership")
	assert.Contains(t, relationshipKinds(fetchedWithRelationships.Relationships), "user-realm-role-mapping")
}

func startIntegrationStack(t *testing.T) {
	t.Helper()
	runIntegrationScript(t, "down.sh", false)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", filepath.Join("mise-tasks", "integration", "up.sh"))
	cmd.Dir = filepath.Join("..")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("integration stack unavailable: %s", string(output))
	}
	t.Cleanup(func() {
		runIntegrationScript(t, "down.sh", true)
	})
}

func runIntegrationScript(t *testing.T, script string, allowFailure bool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", filepath.Join("mise-tasks", "integration", script))
	cmd.Dir = filepath.Join("..")
	output, err := cmd.CombinedOutput()
	if err != nil && !allowFailure {
		t.Skipf("integration script %s failed: %s", script, string(output))
	}
}

func loadEnv(t *testing.T) {
	t.Helper()
	require.NoError(t, godotenv.Overload(filepath.Join("..", ".env")))
}

func bootstrapAdminToken(t *testing.T) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	token, err := auth.New().PasswordToken(ctx, os.Getenv("KEYCLOAK_BASE_URL"), "master", "admin", "admin")
	require.NoError(t, err)
	require.NoError(t, os.Setenv(auth.AccessTokenEnvVar, token.AccessToken))
	require.NoError(t, os.Setenv(auth.RefreshTokenEnvVar, token.RefreshToken))
}

func relationshipKinds(relationships []manifest.RelationshipOperation) []string {
	kinds := make([]string, 0, len(relationships))
	for _, relationship := range relationships {
		kinds = append(kinds, relationship.Kind)
	}
	return kinds
}

func assertRoundTripResourceSubset(t *testing.T, expected []manifest.Resource, actual []manifest.Resource) {
	t.Helper()

	normalizedExpected, _ := manifest.NormalizeRoundTrip(expected, nil)
	normalizedActual, _ := manifest.NormalizeRoundTrip(actual, nil)
	actualByKey := make(map[string]manifest.Resource, len(normalizedActual))
	for _, resource := range normalizedActual {
		actualByKey[roundTripResourceKey(resource)] = resource
	}

	for _, resource := range normalizedExpected {
		actualResource, ok := actualByKey[roundTripResourceKey(resource)]
		require.Truef(t, ok, "missing normalized resource %s", roundTripResourceKey(resource))
		assertMapSubset(t, resource.Data, actualResource.Data)
	}
}

func assertRoundTripRelationshipsEqual(t *testing.T, resources []manifest.Resource, expected []manifest.RelationshipOperation, actual []manifest.RelationshipOperation) {
	t.Helper()

	_, normalizedExpected := manifest.NormalizeRoundTrip(resources, expected)
	_, normalizedActual := manifest.NormalizeRoundTrip(resources, actual)
	require.Equal(t, normalizedExpected, normalizedActual)
}

func roundTripResourceKey(resource manifest.Resource) string {
	name := strings.TrimSpace(resource.Name())
	if name == "" {
		name = strings.TrimSpace(resource.DisplayName())
	}
	return strings.Join([]string{resource.Type, strings.TrimSpace(resource.Realm), name}, "|")
}

func assertMapSubset(t *testing.T, expected map[string]interface{}, actual map[string]interface{}) {
	t.Helper()
	for key, expectedValue := range expected {
		actualValue, ok := actual[key]
		require.Truef(t, ok, "missing key %s", key)
		assertValueSubset(t, expectedValue, actualValue)
	}
}

func assertValueSubset(t *testing.T, expected interface{}, actual interface{}) {
	t.Helper()
	switch expectedTyped := expected.(type) {
	case map[string]interface{}:
		actualTyped, ok := actual.(map[string]interface{})
		require.True(t, ok)
		assertMapSubset(t, expectedTyped, actualTyped)
	case []interface{}:
		actualTyped, ok := actual.([]interface{})
		require.True(t, ok)
		require.GreaterOrEqual(t, len(actualTyped), len(expectedTyped))
		for idx := range expectedTyped {
			assertValueSubset(t, expectedTyped[idx], actualTyped[idx])
		}
	default:
		assert.Equal(t, expected, actual)
	}
}
