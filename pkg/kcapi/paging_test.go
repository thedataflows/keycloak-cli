package kcapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestListAllPages drives ListAll against a fake Keycloak serving 250 users:
// two full pages of 100 and a short page of 50 — the short page ends the
// walk. ListAll always sends first (0, 100, 200), so the server observes
// every offset explicitly; the fake answers each offset with users f..f+100
// capped at 250, mirroring Keycloak's first/max paging.
func TestListAllPages(t *testing.T) {
	var seenFirst []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		first := r.URL.Query().Get("first")
		seenFirst = append(seenFirst, first)
		f := 0
		_, _ = fmt.Sscanf(first, "%d", &f)
		end := f + DefaultPageSize
		if end > 250 {
			end = 250
		}
		users := make([]string, 0, end-f)
		for i := f; i < end; i++ {
			users = append(users, fmt.Sprintf(`{"id":"u-%d"}`, i))
		}
		_, _ = w.Write([]byte("[" + strings.Join(users, ",") + "]"))
	}))
	defer srv.Close()

	c := newInvokeClient(t, srv.URL, realSpecBytes(t))
	out, err := c.ListAll(context.Background(), Call{Resource: "users", Verb: Get, Realm: "demo"})
	require.NoError(t, err)
	assert.Len(t, out, 250)
	assert.Equal(t, []string{"0", "100", "200"}, seenFirst)
}
