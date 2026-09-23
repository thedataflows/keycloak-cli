package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/kcapi"
)

func TestInvokeCmdBuildsCall(t *testing.T) {
	tests := []struct {
		name     string
		opID     string
		realm    string
		params   map[string]string
		body     json.RawMessage
		resource string
		verb     string
		want     kcapi.Call
		wantErr  string
	}{
		{
			name:   "operationId mode maps op, realm, params and body",
			opID:   "getUser",
			realm:  "demo",
			params: map[string]string{"id": "u-1"},
			body:   json.RawMessage(`{"a":1}`),
			want: kcapi.Call{
				Op:     "getUser",
				Realm:  "demo",
				Params: kcapi.P{"id": "u-1"},
				Body:   json.RawMessage(`{"a":1}`),
			},
		},
		{
			name:     "resource+verb mode clears op and sets the verb",
			realm:    "demo",
			params:   map[string]string{"id": "u-1"},
			resource: "users",
			verb:     "GET",
			want: kcapi.Call{
				Resource: "users",
				Verb:     kcapi.Get,
				Realm:    "demo",
				Params:   kcapi.P{"id": "u-1"},
			},
		},
		{
			name:   "no body leaves Body nil",
			opID:   "getUser",
			params: map[string]string{"id": "u-1"},
			want: kcapi.Call{
				Op:     "getUser",
				Params: kcapi.P{"id": "u-1"},
			},
		},
		{
			name:     "operationId and resource together are rejected",
			opID:     "getUser",
			resource: "users",
			verb:     "GET",
			wantErr:  "not both",
		},
		{
			name:     "resource mode requires a verb",
			resource: "users",
			wantErr:  "requires a --verb",
		},
		{
			name:    "no resolution mode is rejected",
			wantErr: "operationId or --resource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := invokeCall(tt.opID, tt.realm, tt.params, tt.body, tt.resource, tt.verb)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInvokeCmdReadBodyArg(t *testing.T) {
	tmp := t.TempDir()
	bodyFile := filepath.Join(tmp, "body.json")
	require.NoError(t, os.WriteFile(bodyFile, []byte(`{"from":"file"}`), 0o644))
	invalidFile := filepath.Join(tmp, "invalid.json")
	require.NoError(t, os.WriteFile(invalidFile, []byte(`{not json`), 0o644))

	tests := []struct {
		name    string
		arg     string
		want    json.RawMessage
		wantErr string
	}{
		{
			name: "inline JSON passes through",
			arg:  `{"inline":true}`,
			want: json.RawMessage(`{"inline":true}`),
		},
		{
			name: "@file reads the file contents",
			arg:  "@" + bodyFile,
			want: json.RawMessage(`{"from":"file"}`),
		},
		{
			name:    "invalid inline JSON is rejected",
			arg:     `not json`,
			wantErr: "body is not valid JSON",
		},
		{
			name:    "invalid file contents are rejected",
			arg:     "@" + invalidFile,
			wantErr: "body is not valid JSON",
		},
		{
			name:    "missing body file is reported",
			arg:     "@" + filepath.Join(tmp, "nope.json"),
			wantErr: "read body file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readBodyArg(tt.arg)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.JSONEq(t, string(tt.want), string(got))
		})
	}
}

func TestInvokeCmdParsesFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(t *testing.T, c InvokeCmd)
	}{
		{
			name: "resource+verb mode with params and body",
			args: []string{
				"invoke",
				"--resource", "users",
				"--verb", "POST",
				"--realm", "demo",
				"--param", "id=u-1",
				"--param", "briefRepresentation=true",
				"--body", `{"username":"alice"}`,
			},
			check: func(t *testing.T, c InvokeCmd) {
				assert.Empty(t, c.OpID)
				assert.Equal(t, "users", c.Resource)
				assert.Equal(t, "POST", c.Verb)
				assert.Equal(t, "demo", c.Realm)
				assert.Equal(t, map[string]string{"id": "u-1", "briefRepresentation": "true"}, c.Param)
				assert.Equal(t, `{"username":"alice"}`, c.Body)
			},
		},
		{
			name: "operationId positional mode",
			args: []string{"invoke", "getUser", "--realm", "demo"},
			check: func(t *testing.T, c InvokeCmd) {
				assert.Equal(t, "getUser", c.OpID)
				assert.Equal(t, "demo", c.Realm)
				assert.Empty(t, c.Resource)
				assert.Equal(t, "GET", c.Verb)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli, kong.Name("keycloak-cli"), kong.Exit(func(int) {}))
			require.NoError(t, err)

			// SpecPath is required; parsing never loads it, so any value works.
			args := append([]string{"--spec-path", "spec.json"}, tt.args...)
			_, err = parser.Parse(args)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.check(t, cli.Invoke)
		})
	}
}
