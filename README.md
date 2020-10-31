# wsjtx-go

Golang binding for the WSJT-X amateur radio software's UDP communication interface. At this early
stage, this only listens and decodes a limited number of message types, but the eventual plan is to
support all message types and two-way communication with WSJT-X.

## Run

This repository is designed as a library but includes a simple driver program to document basic
integration. WSJT-X must be running and generating UDP packets for the driver to pick them up.

From this directory:

```shell script
go run cmd/main.go
```
