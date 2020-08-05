# go-dfms-gateway
HTTP Gateway for DFMS Drive

## Overview
Used for determining Content-Type to render a file in a browser.

## Start using

### Get this repo:

`go get https://github.com/proximax-storage/go-dfms-gateway`

### Quick start

```go
package main

import (
	"log"

	"github.com/proximax-storage/go-dfms-gateway/server"

	apihttp "github.com/proximax-storage/go-xpx-dfms-api-http"
)

func main() {
	// use DFMS API HTTP
	gateway := server.NewGateway(apihttp.NewClientAPI("http://localhost:6366"))
	// start
	log.Fatal(gateway.Start())
}
```

## Test
`go test ./...`

## Build

There is the example - `main.go`.

`go build -o dfms_gateway`

## Run

Run with an address of DFMS API server as argument. 

`./dfms_gateway "api-addr"`

Can ran with flags

### Flags

`-addr` - gateway listening address

`-debug` - enable debug mode

`-cfg` - a path to a custom config file

## Config

The default config path `~/.dfms_gateway/config.json`.

```json
{
	"Name": "DFMS Gateway",
	"Address": ":5000",
	"ApiAddress": "http://localhost:6366",
	"GetOnly": true,
	"LogAllError": false
}
```