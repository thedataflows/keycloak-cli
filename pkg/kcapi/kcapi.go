// Package kcapi is a spec-driven Keycloak client: every operation in the
// supplied OpenAPI spec is discoverable and invocable generically, plus
// relationship-graph queries over Keycloak objects.
package kcapi

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/thedataflows/keycloak-cli/pkg/auth"
)

// defaultTimeout bounds each Keycloak request when Config.Timeout is unset.
const defaultTimeout = 30 * time.Second

// maxSpecBytes caps the size of a spec fetched from SpecSource.URL.
const maxSpecBytes = 32 << 20

// SpecSource describes where New loads the OpenAPI spec from. Exactly one of
// Path, URL, or Raw should be set; Raw wins when several are set.
type SpecSource struct {
	Path string // file path
	URL  string // fetched once at New; Reload refetches
	Raw  []byte // inline bytes
}

// Credentials select which grant the auth service uses: password grant when
// Username/Password are set, client-credentials when ClientSecret is set.
type Credentials struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

// Config configures a Client. Auth, when set, overrides Credentials.
type Config struct {
	BaseURL     string
	Spec        SpecSource
	Credentials Credentials
	Auth        auth.Service // optional; overrides Credentials
	Timeout     time.Duration
	HTTP        *http.Client // optional override
}

// Client is the public entry point of the kcapi library: it owns the loaded
// spec, the runtime transport, and the auth service.
type Client struct {
	spec    *Spec
	runtime *RuntimeClient
	http    *http.Client
	tokens  TokenProvider
	timeout time.Duration
}

// New builds a Client from cfg: it loads the spec, resolves the token
// provider, and wires the runtime transport. No token exchange happens here;
// tokens are acquired lazily per request by the auth service.
func New(cfg Config) (*Client, error) {
	spec, err := loadSpec(cfg.Spec, cfg.HTTP)
	if err != nil {
		return nil, err
	}

	var tokens TokenProvider = cfg.Auth
	if tokens == nil {
		tokens, err = authFromCredentials(cfg.Credentials)
		if err != nil {
			return nil, err
		}
	}

	httpClient := cfg.HTTP
	if httpClient == nil {
		timeout := cfg.Timeout
		if timeout == 0 {
			timeout = defaultTimeout
		}
		httpClient = &http.Client{Timeout: timeout}
	}

	rt, err := NewRuntimeClient(RuntimeConfig{
		BaseURL: cfg.BaseURL,
		Timeout: cfg.Timeout,
		Spec:    spec,
		HTTP:    httpClient,
	}, tokens)
	if err != nil {
		return nil, err
	}

	return &Client{
		spec:    spec,
		runtime: rt,
		http:    httpClient,
		tokens:  tokens,
		timeout: cfg.Timeout,
	}, nil
}

// Spec returns the loaded OpenAPI spec.
func (c *Client) Spec() *Spec {
	if c == nil {
		return nil
	}
	return c.spec
}

// loadSpec dispatches SpecSource on Path/URL/Raw.
func loadSpec(src SpecSource, httpClient *http.Client) (*Spec, error) {
	switch {
	case len(src.Raw) > 0:
		return NewSpecFromBytes(src.Raw)
	case src.Path != "":
		return NewSpec(src.Path)
	case src.URL != "":
		return fetchSpec(src.URL, httpClient)
	default:
		return nil, fmt.Errorf("spec source requires one of Path, URL, or Raw")
	}
}

func fetchSpec(specURL string, httpClient *http.Client) (*Spec, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Get(specURL)
	if err != nil {
		return nil, fmt.Errorf("fetch spec: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch spec: unexpected status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxSpecBytes))
	if err != nil {
		return nil, fmt.Errorf("fetch spec: %w", err)
	}
	return NewSpecFromBytes(body)
}

// authFromCredentials returns the token provider for the given credentials,
// reusing the production pkg/auth service. Password credentials select the
// password grant, a client secret selects the client-credentials grant; the
// grants themselves stay in the auth package's existing methods. With no
// credentials the service resolves tokens from the environment, matching the
// CLI's behavior.
func authFromCredentials(creds Credentials) (TokenProvider, error) {
	hasPassword := creds.Username != "" || creds.Password != ""
	hasClient := creds.ClientSecret != ""
	switch {
	case hasPassword && hasClient:
		return nil, fmt.Errorf("credentials must set either username/password or client secret, not both")
	case hasPassword && (creds.Username == "" || creds.Password == ""):
		return nil, fmt.Errorf("password credentials require both username and password")
	}
	return auth.New(), nil
}
