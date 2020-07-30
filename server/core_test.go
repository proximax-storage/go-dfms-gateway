package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// DFMS node should be ran
func TestNewGateway(t *testing.T) {
	cfg := DefaultConfig()

	s := NewGateway()
	assert.Equal(t, s.address, cfg.Address)
	assert.Equal(t, s.server.Name, cfg.Name)
	assert.Equal(t, s.server.GetOnly, cfg.GetOnly)
	assert.Equal(t, s.server.LogAllErrors, cfg.LogAllError)

	go func() {
		err := s.Start()
		assert.Nil(t, err, err)
	}()

	time.Sleep(5 * time.Second)
	err := s.Stop()
	assert.Nil(t, err, err)
}

func TestNewGatewayWithOptions(t *testing.T) {
	opts := GatewayOptions{
		address:    ":5555",
		apiAddress: "test",
		debug:      true,
	}
	s := NewGateway(
		WithAddress(opts.address),
		Debug(true),
	)
	assert.Equal(t, opts.address, s.address)
	assert.Equal(t, opts.debug, s.server.LogAllErrors)

	go func() {
		err := s.Start()
		assert.Nil(t, err, err)
	}()

	time.Sleep(5 * time.Second)
	err := s.Stop()
	assert.Nil(t, err, err)
}
