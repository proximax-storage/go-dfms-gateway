package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"path/filepath"
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

func TestResolvePath(t *testing.T) {
	cases := []struct {
		val    string
		expect func() string
	}{
		{
			"~/path",
			func() string {
				home, err := os.UserHomeDir()
				if err != nil {
					log.Fatalf("Cannot get user home dir: ", err)
				}

				return path.Join(home, "path")
			},
		},
		{
			"~path",
			func() string {
				return "~path"
			},
		},
		{
			"/path",
			func() string {
				return "/path"
			},
		},
		{
			"path",
			func() string {
				return "path"
			},
		},
		{
			"path~/",
			func() string {
				return "path~/"
			},
		},
	}

	for _, v := range cases {
		r := resolvePath(v.val)
		assert.Equal(t, v.expect(), r)
	}
}

func TestMatch(t *testing.T) {
	cases := []struct {
		val    []byte
		expect bool
	}{
		{
			[]byte("baegaajaiaqjcb5e3ur46es6qckg5aaayqcnqxxtmffp37lfkitfjlmenmlopvt2h"),
			true,
		},
		{
			[]byte("baegaajaiaqjcb5e3ur46es6qckg5aaayqcnqxxtmffp37lf"),
			false,
		},
		{
			[]byte("baegaajaiaqjcb5e3ur46es6qckg5aaayqcnqxxtmffp37lfkitfjlmenmlopvt2hsdsdsds"),
			true,
		},
		{
			[]byte("afdsdsdbaegaajaiaqjcb5e3ur46es6qckg5aaayqcnqxxtmffp37lfkitfjlmenmlopvt2h"),
			true,
		},
		{
			[]byte("baegaajaiaqjcb5e3ur46es6qckg5aaa/yqcnqxxtmffp37lfkitfjlmenmlopvt2"),
			false,
		},
	}

	for _, v := range cases {
		b := match(cidPattern, v.val)
		assert.Equal(t, v.expect, b)
	}
}

func TestConfig_Save_Load(t *testing.T) {
	cfg := DefaultConfig()

	p, err := filepath.Abs("test.json")
	assert.Nil(t, err, err)

	err = saveConfig(cfg, p)
	assert.Nil(t, err, err)

	loadedCfg, err := loadConfig(p)
	assert.Nil(t, err, err)
	require.NotNil(t, loadedCfg)
	assert.Equal(t, *cfg, *loadedCfg)

	err = os.Remove(p)
	assert.Nil(t, err, err)
}
