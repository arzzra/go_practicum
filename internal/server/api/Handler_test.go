package api

import (
	"github.com/arzzra/go_practicum/internal/server"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeHandler(t *testing.T) {
	_, err := MakeHandler(nil)
	assert.NotNil(t, err)

	srv, err := MakeHandler(&server.Server{})
	assert.Nil(t, err)
	assert.NotNil(t, srv)
}
