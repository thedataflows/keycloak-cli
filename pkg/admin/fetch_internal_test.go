package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildQueryParamsFullRepresentation(t *testing.T) {
	tests := []struct {
		name  string
		query FetchQuery
		want  map[string]string
	}{
		{
			name:  "no options emits no params",
			query: FetchQuery{},
			want:  nil,
		},
		{
			name:  "full representation alone emits briefRepresentation false",
			query: FetchQuery{FullRepresentation: true},
			want:  map[string]string{"briefRepresentation": "false"},
		},
		{
			name:  "brief default omits the key entirely",
			query: FetchQuery{Search: "acme", Max: 10},
			want:  map[string]string{"search": "acme", "max": "10"},
		},
		{
			name:  "full representation combines with search",
			query: FetchQuery{Search: "acme", FullRepresentation: true},
			want:  map[string]string{"search": "acme", "briefRepresentation": "false"},
		},
		{
			name:  "full representation combines with max",
			query: FetchQuery{Max: 5, FullRepresentation: true},
			want:  map[string]string{"max": "5", "briefRepresentation": "false"},
		},
		{
			name:  "full representation combines with search and max",
			query: FetchQuery{Search: "acme", Max: 5, FullRepresentation: true},
			want:  map[string]string{"search": "acme", "max": "5", "briefRepresentation": "false"},
		},
		{
			name:  "full representation combines with exact match",
			query: FetchQuery{Search: "acme", ExactMatch: true, FullRepresentation: true},
			want:  map[string]string{"search": "acme", "exact": "true", "briefRepresentation": "false"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildQueryParams(tt.query)

			if tt.want == nil {
				assert.Nil(t, got)
				return
			}

			require.Len(t, got, 1)
			assert.Equal(t, tt.want, got[0])
		})
	}
}
