openapi3-gen
============

A tool to generate OpenAPI 3 specification using annotations embedded in Go
source code.

Installation
------------
`openapi3-gen` uses Go modules.

You can install it from source:
```
go get -u github.com/skygeario/openapi3-gen/cmd/openapi3-gen
```

Usage
-----
```
usage: openapi3-gen [flags] <patterns...>
  -d string
  -dir string
        project base directory (default to working directory)
  -o string
  -output string
        output OpenAPI specification file
```

For example, if source code is placed in `/project/cmd/` and `/project/pkg/`:
```
openapi3-gen -d /project -o /project/docs/api.yaml ./cmd/... ./pkg/...
```
OpenAPI 3 specification would be saved to `/project/docs/api.yaml`.

Example usages can be found in [`/examples`](./examples).

License
-------
```
Copyright (c) 2019-present, Oursky Ltd.
All rights reserved.

This source code is licensed under the Apache License version 2.0
found in the LICENSE file in the root directory of this source tree.
An additional grant of patent rights can be found in the PATENTS
file in the same directory.
```
