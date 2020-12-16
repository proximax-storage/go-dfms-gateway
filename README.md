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

`cd cmd && go build -o dfms_gateway`

## Run

Run with an address of DFMS API server as argument. 

`./dfms_gateway "api-addr"`

Can ran with flags

### Flags

- `-addr` - gateway listening address
- `-debug` - enable debug mode
- `-cors` - enable cors
- `-methods` - List of allowed CORs methods separated by commas.
- `-headers` - List of allowed CORs headers separated by commas.
- `-origins` - List of allowed CORs origins separated by commas.

## Config

The default config path `~/.dfms-client_gateway/config.json`.

```json
{
  "Name": "DFMS Gateway",
  "Address": ":5000",
  "GetOnly": false,
  "LogAllError": false,
  "CORs": {
    "Enable": true,
    "AllowedMethods": [
      "*"
    ],
    "AllowedHeaders": [
      "*"
    ],
    "AllowedOrigins": [
      "*"
    ]
  }
}
```