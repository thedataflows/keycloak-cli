package kcapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang.org/x/oauth2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// syntheticSpec exists because the vendored 26.6.2 spec carries zero
// operationId fields: Op-mode resolution needs operationIds to be
// deterministic, so every Op-mode Invoke test runs against this minimal
// inline spec (one single-user GET, one collection POST) instead.
const syntheticSpec = `{
  "openapi": "3.0.3",
  "info": {"title": "kcapi invoke test", "version": "0.0.0"},
  "paths": {
    "/admin/realms/{realm}/users/{id}": {
      "get": {
        "operationId": "getUser",
        "tags": ["Users"],
        "summary": "Get representation of the user",
        "parameters": [
          {"name": "realm", "in": "path", "required": true, "schema": {"type": "string"}},
          {"name": "id", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "responses": {
          "200": {"description": "OK", "content": {"application/json": {"schema": {"type": "object"}}}}
        }
      }
    },
    "/admin/realms/{realm}/users": {
      "post": {
        "operationId": "createUser",
        "tags": ["Users"],
        "summary": "Create a user",
        "parameters": [
          {"name": "realm", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["username"],
                "properties": {"username": {"type": "string"}}
              }
            }
          }
        },
        "responses": {"201": {"description": "Created"}}
      }
    }
  }
}`

// duplicateOpSpec maps one operationId onto two distinct paths: resolving it
// must be a validation error, never a silent pick.
const duplicateOpSpec = `{
  "openapi": "3.0.3",
  "info": {"title": "kcapi duplicate op test", "version": "0.0.0"},
  "paths": {
    "/admin/realms/{realm}/users/{id}": {
      "get": {
        "operationId": "getUser",
        "parameters": [
          {"name": "realm", "in": "path", "required": true, "schema": {"type": "string"}},
          {"name": "id", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "responses": {"200": {"description": "OK"}}
      }
    },
    "/admin/realms/{realm}/clients/{id}": {
      "get": {
        "operationId": "getUser",
        "parameters": [
          {"name": "realm", "in": "path", "required": true, "schema": {"type": "string"}},
          {"name": "id", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "responses": {"200": {"description": "OK"}}
      }
    }
  }
}`

// staticTokenProvider is a fixed-token auth.Service: AccessToken always
// returns the configured token, so Invoke tests need no token endpoint.
type staticTokenProvider string

func (p staticTokenProvider) AccessToken(_ context.Context, _, _, _ string) (string, error) {
	return string(p), nil
}

func (p staticTokenProvider) PasswordToken(context.Context, string, string, string, string) (oauth2.Token, error) {
	return oauth2.Token{}, nil
}

func (p staticTokenProvider) ClientCredentialsToken(context.Context, string, string, string, string) (oauth2.Token, error) {
	return oauth2.Token{}, nil
}

func (p staticTokenProvider) SetEnvToken(string, string, string) error { return nil }

// realSpecBytes loads the vendored Keycloak spec for resource+verb tests.
func realSpecBytes(t *testing.T) []byte {
	t.Helper()
	raw, err := os.ReadFile("../../keycloak-oapi/26.7.4.spec.json")
	if err != nil {
		t.Skipf("repo spec not available: %v", err)
	}
	return raw
}

func newInvokeClient(t *testing.T, baseURL string, spec []byte) *Client {
	t.Helper()
	c, err := New(Config{
		BaseURL: baseURL,
		Spec:    SpecSource{Raw: spec},
		Auth:    staticTokenProvider("test-token"),
	})
	require.NoError(t, err)
	return c
}

// capturedRequest records what actually hit the fake Keycloak server.
type capturedRequest struct {
	Method        string
	Path          string
	RawQuery      string
	Authorization string
	Accept        string
	ContentType   string
	Body          []byte
}

// newInvokeServer serves one canned response and captures the request; spec
// selects the synthetic inline spec or the real vendored one.
func newInvokeServer(t *testing.T, spec []byte, status int, responseBody string) (*Client, *capturedRequest) {
	t.Helper()
	captured := &capturedRequest{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		*captured = capturedRequest{
			Method:        r.Method,
			Path:          r.URL.Path,
			RawQuery:      r.URL.RawQuery,
			Authorization: r.Header.Get("Authorization"),
			Accept:        r.Header.Get("Accept"),
			ContentType:   r.Header.Get("Content-Type"),
			Body:          body,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(responseBody))
	}))
	t.Cleanup(srv.Close)
	return newInvokeClient(t, srv.URL, spec), captured
}

func TestInvokeByOperationID(t *testing.T) {
	t.Run("resolves operationId and fills path params", func(t *testing.T) {
		c, captured := newInvokeServer(t, []byte(syntheticSpec), http.StatusOK, `{"id":"u-1","username":"alice"}`)

		out, err := c.Invoke(context.Background(), Call{Op: "getUser", Realm: "demo", Params: P{"id": "u-1"}})
		require.NoError(t, err)

		assert.Equal(t, http.MethodGet, captured.Method)
		assert.Equal(t, "/admin/realms/demo/users/u-1", captured.Path)
		assert.Equal(t, "Bearer test-token", captured.Authorization)
		assert.Equal(t, "application/json", captured.Accept)

		var user map[string]interface{}
		require.NoError(t, json.Unmarshal(out, &user))
		assert.Equal(t, "alice", user["username"])
	})

	t.Run("Params override Realm for the realm placeholder", func(t *testing.T) {
		c, captured := newInvokeServer(t, []byte(syntheticSpec), http.StatusOK, `{}`)

		_, err := c.Invoke(context.Background(), Call{Op: "getUser", Realm: "demo", Params: P{"id": "u-1", "realm": "master"}})
		require.NoError(t, err)

		assert.Equal(t, "/admin/realms/master/users/u-1", captured.Path)
	})
}

func TestInvokeValidationErrors(t *testing.T) {
	// Unreachable base URL: any of these cases that tried to send a request
	// would surface as a network error, not ErrValidation.
	const unreachable = "http://kcapi-invalid.invalid"

	t.Run("unknown operationId", func(t *testing.T) {
		c := newInvokeClient(t, unreachable, []byte(syntheticSpec))

		_, err := c.Invoke(context.Background(), Call{Op: "no-such-op"})
		assert.ErrorIs(t, err, ErrValidation)
	})

	t.Run("missing required path param", func(t *testing.T) {
		c := newInvokeClient(t, unreachable, []byte(syntheticSpec))

		_, err := c.Invoke(context.Background(), Call{Op: "getUser", Realm: "demo"})
		assert.ErrorIs(t, err, ErrValidation)
	})

	t.Run("op and resource modes are mutually exclusive", func(t *testing.T) {
		c := newInvokeClient(t, unreachable, []byte(syntheticSpec))

		_, err := c.Invoke(context.Background(), Call{Op: "getUser", Resource: "users", Verb: Get})
		assert.ErrorIs(t, err, ErrValidation)
	})

	t.Run("neither resolution mode given", func(t *testing.T) {
		c := newInvokeClient(t, unreachable, []byte(syntheticSpec))

		_, err := c.Invoke(context.Background(), Call{})
		assert.ErrorIs(t, err, ErrValidation)
	})

	t.Run("body forbidden on operation without request schema", func(t *testing.T) {
		c := newInvokeClient(t, unreachable, []byte(syntheticSpec))

		_, err := c.Invoke(context.Background(), Call{Op: "getUser", Realm: "demo", Params: P{"id": "u-1"}, Body: map[string]any{"username": "x"}})
		assert.ErrorIs(t, err, ErrValidation)
	})

	t.Run("duplicate operationId is a validation error", func(t *testing.T) {
		c := newInvokeClient(t, unreachable, []byte(duplicateOpSpec))

		_, err := c.Invoke(context.Background(), Call{Op: "getUser", Realm: "demo", Params: P{"id": "x"}})
		assert.ErrorIs(t, err, ErrValidation)
	})
}

func TestInvokeByResourceAndVerb(t *testing.T) {
	t.Run("GET single user via the real user-id placeholder", func(t *testing.T) {
		c, captured := newInvokeServer(t, realSpecBytes(t), http.StatusOK, `{"id":"u-1","username":"alice"}`)

		out, err := c.Invoke(context.Background(), Call{Resource: "users", Verb: Get, Realm: "demo", Params: P{"user-id": "u-1"}})
		require.NoError(t, err)

		assert.Equal(t, http.MethodGet, captured.Method)
		assert.Equal(t, "/admin/realms/demo/users/u-1", captured.Path)

		var user map[string]interface{}
		require.NoError(t, json.Unmarshal(out, &user))
		assert.Equal(t, "alice", user["username"])
	})

	t.Run("POST collection create forwards the JSON body", func(t *testing.T) {
		c, captured := newInvokeServer(t, realSpecBytes(t), http.StatusCreated, `{"id":"u-2"}`)

		out, err := c.Invoke(context.Background(), Call{Resource: "users", Verb: Post, Realm: "demo", Body: map[string]any{"username": "bob"}})
		require.NoError(t, err)

		assert.Equal(t, http.MethodPost, captured.Method)
		assert.Equal(t, "/admin/realms/demo/users", captured.Path)
		assert.Equal(t, "application/json", captured.ContentType)

		var sent map[string]interface{}
		require.NoError(t, json.Unmarshal(captured.Body, &sent))
		assert.Equal(t, "bob", sent["username"])

		var created map[string]interface{}
		require.NoError(t, json.Unmarshal(out, &created))
		assert.Equal(t, "u-2", created["id"])
	})

	t.Run("realm-only GET resolves to the shortest collection path", func(t *testing.T) {
		c, captured := newInvokeServer(t, realSpecBytes(t), http.StatusOK, `[{"username":"alice"}]`)

		_, err := c.Invoke(context.Background(), Call{Resource: "users", Verb: Get, Realm: "demo"})
		require.NoError(t, err)

		// Shortest among the fillable users-GET paths (beats /users/count,
		// /users/profile, /users/profile/metadata), and the only candidate
		// whose every placeholder Realm alone fills.
		assert.Equal(t, "/admin/realms/demo/users", captured.Path)
	})

	t.Run("non-path params are forwarded as query params", func(t *testing.T) {
		c, captured := newInvokeServer(t, realSpecBytes(t), http.StatusOK, `[]`)

		_, err := c.Invoke(context.Background(), Call{Resource: "users", Verb: Get, Realm: "demo", Params: P{"first": "5"}})
		require.NoError(t, err)

		assert.Equal(t, "/admin/realms/demo/users", captured.Path)
		assert.Equal(t, "first=5", captured.RawQuery)
	})
}

func TestInvokeValidationErrorsByResourceAndVerb(t *testing.T) {
	const unreachable = "http://kcapi-invalid.invalid"

	t.Run("empty value for a required real placeholder", func(t *testing.T) {
		c := newInvokeClient(t, unreachable, realSpecBytes(t))

		_, err := c.Invoke(context.Background(), Call{Resource: "users", Verb: Get, Realm: "demo", Params: P{"user-id": ""}})
		assert.ErrorIs(t, err, ErrValidation)
	})

	t.Run("body forbidden on GET collection", func(t *testing.T) {
		c := newInvokeClient(t, unreachable, realSpecBytes(t))

		_, err := c.Invoke(context.Background(), Call{Resource: "users", Verb: Get, Realm: "demo", Body: map[string]any{"username": "x"}})
		assert.ErrorIs(t, err, ErrValidation)
	})

	t.Run("missing verb", func(t *testing.T) {
		c := newInvokeClient(t, unreachable, realSpecBytes(t))

		_, err := c.Invoke(context.Background(), Call{Resource: "users", Realm: "demo"})
		assert.ErrorIs(t, err, ErrValidation)
	})
}

func TestInvokeEmpty204Body(t *testing.T) {
	c, captured := newInvokeServer(t, realSpecBytes(t), http.StatusNoContent, "")

	out, err := c.Invoke(context.Background(), Call{Resource: "users", Verb: Put, Realm: "demo", Params: P{"user-id": "u-1"}, Body: map[string]any{"firstName": "Alice"}})
	require.NoError(t, err)

	assert.Nil(t, out, "a 204 with an empty body is a valid nil result")
	assert.Equal(t, http.MethodPut, captured.Method)
	assert.Equal(t, "/admin/realms/demo/users/u-1", captured.Path)
	assert.Contains(t, string(captured.Body), "firstName")
}

// A json.RawMessage body must validate as its decoded JSON value, not as a
// byte slice: cmd/invoke.go hands RawMessage bodies straight from --body.
func TestValidateCallAcceptsRawMessageBody(t *testing.T) {
	spec, err := NewSpecFromBytes(realSpecBytes(t))
	require.NoError(t, err)

	call := Call{
		Resource: "users",
		Verb:     Post,
		Realm:    "master",
		Body:     json.RawMessage(`{"username":"bob"}`),
	}
	require.NoError(t, ValidateCall(spec, "/admin/realms/{realm}/users", http.MethodPost, call),
		"a valid RawMessage body must pass validation")

	bad := call
	bad.Body = json.RawMessage(`{not json`)
	err = ValidateCall(spec, "/admin/realms/{realm}/users", http.MethodPost, bad)
	require.Error(t, err, "a malformed RawMessage body must fail validation")
}
