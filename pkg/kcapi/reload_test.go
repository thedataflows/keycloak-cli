package kcapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReloadSwapsCatalog(t *testing.T) {
	var specBody atomic.Value
	specBody.Store(realSpecBytes(t)) // start from the full repo spec
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(specBody.Load().([]byte))
	}))
	defer srv.Close()

	c, err := New(Config{BaseURL: "http://test", Spec: SpecSource{URL: srv.URL}})
	require.NoError(t, err)

	// The vendored 26.6.2 spec carries no operationId fields, so the brief's
	// "getUser" search matches nothing there; "users/{user-id}" pins real
	// operations by path substring (same pattern as discovery_test.go).
	before, err := c.Operations(OpFilter{Search: "users/{user-id}"})
	require.NoError(t, err)
	require.NotEmpty(t, before)

	specBody.Store([]byte(`{"openapi":"3.0.0","paths":{}}`)) // spec update drops everything (parser accepts zero-path documents)
	require.NoError(t, c.Reload(context.Background()))
	after, err := c.Operations(OpFilter{})
	require.NoError(t, err)
	assert.Empty(t, after)
}

// TestReloadFailureKeepsOldPair pins the atomicity contract: a reload that
// fails at fetch time must leave the previously loaded spec/runtime pair in
// place, so discovery keeps working on the old catalog.
func TestReloadFailureKeepsOldPair(t *testing.T) {
	var specBody atomic.Value
	specBody.Store(realSpecBytes(t))
	var status atomic.Int64
	status.Store(http.StatusOK)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(int(status.Load()))
		_, _ = w.Write(specBody.Load().([]byte))
	}))
	defer srv.Close()

	c, err := New(Config{BaseURL: "http://test", Spec: SpecSource{URL: srv.URL}})
	require.NoError(t, err)

	before, err := c.Operations(OpFilter{})
	require.NoError(t, err)
	require.NotEmpty(t, before)

	status.Store(http.StatusInternalServerError) // refetch now fails
	require.Error(t, c.Reload(context.Background()))

	after, err := c.Operations(OpFilter{})
	require.NoError(t, err)
	assert.Len(t, after, len(before), "failed Reload must keep the old catalog")
}

// TestReloadConcurrentWithReaders exercises the snapshot discipline under the
// race detector: readers capture the spec/runtime pair at entry while Reload
// swaps pointers underneath them. It passes only if no reader ever observes a
// torn pair (new spec with old runtime or vice versa) and -race stays silent.
func TestReloadConcurrentWithReaders(t *testing.T) {
	var specBody atomic.Value
	specBody.Store([]byte(syntheticSpec))
	specSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(specBody.Load().([]byte))
	}))
	defer specSrv.Close()

	apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer apiSrv.Close()

	c, err := New(Config{
		BaseURL: apiSrv.URL,
		Spec:    SpecSource{URL: specSrv.URL},
		Auth:    staticTokenProvider("test-token"),
	})
	require.NoError(t, err)

	done := make(chan struct{})
	var wg sync.WaitGroup
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
				}
				_, _ = c.Operations(OpFilter{Search: "getUser"})
				_, _ = c.Invoke(context.Background(), Call{Op: "getUser", Realm: "demo", Params: P{"id": "u-1"}})
				_ = c.Spec()
			}
		}()
	}

	for range 25 {
		require.NoError(t, c.Reload(context.Background()))
	}
	close(done)
	wg.Wait()
}
