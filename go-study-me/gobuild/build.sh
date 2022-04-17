#!/usr/bin/env bash

buildtime="$(date -u '+%Y-%m-%d %I:%M:%S%p')"
BRANCH=`git rev-parse --abbrev-ref HEAD`
COMMIT=`git rev-parse --short HEAD`
GOVERSION=`go version`
GOLDFLAGS="-s -w -X 'main.buildtime=$buildtime' -X 'main.branch=$BRANCH' -X 'main.commit=$COMMIT' -X 'main.goversion=$GOVERSION'"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$GOLDFLAGS" -o hello.exe main.go