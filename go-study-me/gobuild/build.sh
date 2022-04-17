#!/usr/bin/env bash

buildDate="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
gitVersion=`git rev-parse --abbrev-ref HEAD`
gitCommit=`git rev-parse --short HEAD`
goVersion="$(go env GOVERSION)"
Platform="$(go env GOOS)/$(go env GOARCH)"
GOLDFLAGS="-s -w \
    -X 'github.com/panlq/gobuild/version.buildDate=$buildDate' \
    -X 'github.com/panlq/gobuild/version.gitCommit=$gitCommit' \
    -X 'github.com/panlq/gobuild/version.gitVersion=$gitVersion' \
    -X 'github.com/panlq/gobuild/version.goVersion=$goVersion' \
    -X 'github.com/panlq/gobuild/version.platform=$Platform' "

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$GOLDFLAGS" -o hello main.go