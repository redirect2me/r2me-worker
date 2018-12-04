#!/bin/bash


go build \
	-ldflags "-X main.COMMIT=$(git rev-parse --short HEAD) -X main.LASTMOD=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
	r2server.go status.go
