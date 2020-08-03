# go-dfms-gateway
HTTP Gateway for DFMS Drive

## Overview
Used for determining Content-Type to render a file in a browser.

## Build

```go build -o dfms_gateway```

## Run

```./dfms_gateway```

Can ran with flags

### Flags

`-addr` - gateway listening address

`-api-addr` - API address of a DFMS node

`-debug` - enable debug mode

`-cfg` - a path to a custom config file

## Config

The default config path `~/.dfms/gateway_cfg.json`.

```json
{
	"Name": "DFMS Gateway",
	"Address": ":5000",
	"ApiAddress": "http://localhost:6366",
	"GetOnly": true,
	"LogAllError": false
}
```