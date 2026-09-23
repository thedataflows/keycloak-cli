package kcapi

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClassifyErrorKinds(t *testing.T) {
	tests := []struct {
		status int
		want   ErrorKind
		sent   error
	}{
		{http.StatusNotFound, KindNotFound, ErrNotFound},
		{http.StatusConflict, KindConflict, ErrConflict},
		{http.StatusUnauthorized, KindAuth, ErrAuth},
		{http.StatusBadRequest, KindValidation, nil},
		{http.StatusInternalServerError, KindServer, nil},
	}
	for _, tt := range tests {
		err := classifyError(errors.New("boom"), tt.status, "invoke getUser")
		assert.Equal(t, tt.want, err.Kind)
		assert.Equal(t, "invoke getUser", err.Op)
		assert.Equal(t, tt.status, err.Status)
		if tt.sent != nil {
			assert.ErrorIs(t, err, tt.sent)
		}
	}
}

func TestErrorIsSentinel(t *testing.T) {
	err := &Error{Kind: KindNotFound, Op: "op"}
	assert.ErrorIs(t, err, ErrNotFound)
	assert.NotErrorIs(t, err, ErrConflict)
}
