// Package kcapi is a spec-driven Keycloak client: every operation in the
// supplied OpenAPI spec is discoverable and invocable generically, plus
// relationship-graph queries over Keycloak objects.
package kcapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
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

// Credentials declare which grant shape you intend: a username/password pair
// or a client secret. They are validated for shape only — contradictory or
// partial input is rejected at New — and their values are not consumed:
// tokens resolve from the environment, or from Config.Auth when it is set.
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
// spec, the runtime transport, and the auth service. The spec/runtime pair is
// guarded by mu: readers capture it via snapshot and keep their captured
// values for the whole call, so Reload can swap the pair atomically without
// disturbing in-flight operations.
type Client struct {
	mu      sync.RWMutex // guards spec and runtime
	spec    *Spec
	runtime *RuntimeClient
	http    *http.Client
	tokens  TokenProvider
	// timeout is the EFFECTIVE per-request timeout: cfg.Timeout when set,
	// defaultTimeout otherwise. It is resolved at New (not stored raw) so
	// Reload rebuilds the runtime with the same timeout the client was
	// constructed with, even when Config.Timeout was left zero.
	timeout time.Duration
	// specSource records where the spec was loaded from; Reload refetches it.
	specSource SpecSource
}

// New builds a Client from cfg: it loads the spec, resolves the token
// provider, and wires the runtime transport. No token exchange happens here;
// tokens are acquired lazily per request by the auth service.
func New(cfg Config) (*Client, error) {
	spec, err := loadSpec(context.Background(), cfg.Spec, cfg.HTTP)
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

	// Resolve the effective timeout once: the raw cfg.Timeout may be zero,
	// which would leave a rebuilt runtime without any request deadline.
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	httpClient := cfg.HTTP
	if httpClient == nil {
		httpClient = &http.Client{Timeout: timeout}
	}

	rt, err := NewRuntimeClient(RuntimeConfig{
		BaseURL: cfg.BaseURL,
		Timeout: timeout,
		Spec:    spec,
		HTTP:    httpClient,
	}, tokens)
	if err != nil {
		return nil, err
	}

	return &Client{
		spec:       spec,
		runtime:    rt,
		http:       httpClient,
		tokens:     tokens,
		timeout:    timeout,
		specSource: cfg.Spec,
	}, nil
}

// snapshot returns the current spec/runtime pair under the read lock. Callers
// must use the returned values for the entire call and never re-read the
// client fields: an in-flight operation keeps the pair it captured even when
// Reload swaps in a new one.
func (c *Client) snapshot() (*Spec, *RuntimeClient) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.spec, c.runtime
}

// Spec returns the loaded OpenAPI spec.
func (c *Client) Spec() *Spec {
	if c == nil {
		return nil
	}
	spec, _ := c.snapshot()
	return spec
}

// Reload refetches the spec from the client's configured source (URL or Path;
// Raw reloads the same bytes), parses it, builds a fresh RuntimeClient, and
// swaps both under one write lock. The new pair is built completely before
// the lock, so a failure at any step — fetch, parse, or runtime build —
// returns early and leaves the old pair untouched. In-flight calls keep the
// old pair: they captured it under the read lock and never re-read the
// client fields.
func (c *Client) Reload(ctx context.Context) error {
	_, oldRuntime := c.snapshot()

	spec, err := loadSpec(ctx, c.specSource, c.http)
	if err != nil {
		return err
	}
	rt, err := NewRuntimeClient(RuntimeConfig{
		BaseURL: oldRuntime.BaseURL(),
		Timeout: c.timeout,
		Spec:    spec,
		HTTP:    c.http,
	}, c.tokens)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.spec, c.runtime = spec, rt
	return nil
}

// loadSpec dispatches SpecSource on Path/URL/Raw.
func loadSpec(ctx context.Context, src SpecSource, httpClient *http.Client) (*Spec, error) {
	switch {
	case len(src.Raw) > 0:
		return NewSpecFromBytes(src.Raw)
	case src.Path != "":
		return NewSpec(src.Path)
	case src.URL != "":
		return fetchSpec(ctx, src.URL, httpClient)
	default:
		return nil, fmt.Errorf("spec source requires one of Path, URL, or Raw")
	}
}

func fetchSpec(ctx context.Context, specURL string, httpClient *http.Client) (*Spec, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, specURL, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch spec: %w", err)
	}
	resp, err := httpClient.Do(req)
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

// authFromCredentials validates the shape of the given credentials —
// username/password and a client secret are mutually exclusive, and a
// password pair must be complete — and returns the production pkg/auth
// service. No grant is selected here and the credential values are not
// consumed: auth.New() resolves tokens from the environment, matching the
// CLI's behavior; Config.Auth overrides it wholesale.
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
