[![PkgGoDev](https://pkg.go.dev/badge/github.com/k0swe/wsjtx-go)](https://pkg.go.dev/github.com/k0swe/wsjtx-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/k0swe/wsjtx-go)](https://goreportcard.com/report/github.com/k0swe/wsjtx-go)
[![Test](https://github.com/k0swe/wsjtx-go/workflows/Test/badge.svg?branch=main)](https://github.com/k0swe/wsjtx-go/actions/workflows/test.yml?query=branch%3Amain)

# wsjtx-go

Golang binding for the WSJT-X amateur radio software's UDP communication interface. This library
supports receiving and sending all WSJT-X message types up through WSJT-X v2.3.0.

This is meant to be a fairly thin binding API, so familiarity with WSJT-X's
[`NetworkMessage.hpp`](https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp)
is recommended.

## Run

This repository is designed as a library but includes a simple driver program to document basic
integration. WSJT-X must be running and generating UDP packets for the driver to pick them up.

From this directory:

```shell script
go run cmd/main.go
```
