package kcapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceRoundTrip(t *testing.T) {
	r := Resource{Type: "users", Data: map[string]interface{}{"username": "alice"}}
	assert.Equal(t, "users", r.Type)
	assert.Equal(t, "alice", r.Data["username"])
}
