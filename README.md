# gojsonrpc
[![GoDoc Reference](https://godoc.org/github.com/asib/gojsonrpc?status.svg)](https://godoc.org/github.com/asib/gojsonrpc)
[![Build Status](https://travis-ci.org/asib/gojsonrpc.svg?branch=master)](https://travis-ci.org/asib/gojsonrpc)
[![Coverage Status](https://coveralls.io/repos/github/asib/gojsonrpc/badge.svg?branch=master)](https://coveralls.io/github/asib/gojsonrpc?branch=master)
[![GitHub tag](https://img.shields.io/github/tag/asib/gojsonrpc.svg?maxAge=2592000)]()
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg?style=plastic)](https://raw.githubusercontent.com/asib/gojsonrpc/master/LICENSE)

A library for working with JSON-RPC APIs in Go.

## Installation

To install, run

```bash
$ go get -u github.com/asib/gojsonrpc
```

## Example

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/asib/gojsonrpc"
)

func handleRequest(req *gojsonrpc.Request) (resp *gojsonrpc.Response) {
	if req.Method() == "login" {
		params := req.Params().(map[string]interface{})

		if params["user"] == "asib" && params["password"] == "pass123" {
			resp = gojsonrpc.MakeResponseWithResult(map[string]bool{"success": true}, req.ID())
		} else {
			resp = gojsonrpc.MakeResponseWithResult(map[string]bool{"success": false}, req.ID())
		}
	}

	return
}

func main() {
	raw := `{"jsonrpc":"2.0", "method":"login", "params":{"user":"asib", "password":"pass123"}, "id":1}`
	msg, err := gojsonrpc.ParseIncoming(raw)
	if err != nil {
		panic(err)
	}

	switch msg.(type) {
	case *gojsonrpc.Notification:
		// do stuff here
	case *gojsonrpc.Request:
		rawResp, err := json.Marshal(handleRequest(msg.(*gojsonrpc.Request)))
		if err != nil {
			panic(err)
		} else {
			fmt.Println(string(rawResp))
		}
	case *gojsonrpc.Response:
		// do stuff here
	default:
	  // handle error case
	}
}
```

## Documentation

Visit the [godoc](https://godoc.org/github.com/asib/gojsonrpc) page.
