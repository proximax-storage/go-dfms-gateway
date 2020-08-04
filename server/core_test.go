package server

import (
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apihttp "github.com/proximax-storage/go-xpx-dfms-api-http"
)

var clientApi = apihttp.NewClientAPI("http://localhost:6366")

// DFMS node should be ran
func TestNewGateway(t *testing.T) {
	cfg := defaultConfig()

	gateway := NewGateway(clientApi)
	assert.Equal(t, gateway.address, cfg.Address)
	assert.Equal(t, gateway.server.Name, cfg.Name)
	assert.Equal(t, gateway.server.GetOnly, cfg.GetOnly)
	assert.Equal(t, gateway.server.LogAllErrors, cfg.LogAllError)

	go func() {
		err := gateway.Start()
		assert.Nil(t, err, err)
	}()

	time.Sleep(5 * time.Second)
	err := gateway.Stop()
	assert.Nil(t, err, err)
}

func TestNewGatewayWithOptions(t *testing.T) {
	opts := gatewayOptions{
		address: ":5555",
		debug:   true,
	}
	gateway := NewGateway(
		clientApi,
		WithAddress(opts.address),
		Debug(true),
	)
	assert.Equal(t, opts.address, gateway.address)
	assert.Equal(t, opts.debug, gateway.server.LogAllErrors)

	go func() {
		err := gateway.Start()
		assert.Nil(t, err, err)
	}()

	time.Sleep(5 * time.Second)
	err := gateway.Stop()
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
	cfg := defaultConfig()

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
