#!/usr/bin/env bash
set -eux;
buildDate="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
gitCommit="$(git rev-parse --short HEAD)"
goVersion="$(go version)"
GOLDFLAGS="-s -w \
    -X 'github.com/panlq/gobuild/version.buildDate=$buildDate' \
    -X 'github.com/panlq/gobuild/version.gitCommit=$gitCommit' \
    -X 'github.com/panlq/gobuild/version.goVersion=$goVersion'"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$GOLDFLAGS" -o main main.go