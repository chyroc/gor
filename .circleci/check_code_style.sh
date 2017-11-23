#!/usr/bin/env bash

set -e

FMT_STATUS=`gofmt -d $(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./_project/*")`
if ["$FMT_STATUS" == ""]; then exit 0; else exit 1; fi

sdaasdf asf
ls
echo $(go list ./... | grep -v /vendor/)

LINT_STATUS=`golint $(go list ./... | grep -v /vendor/)`
if ["$LINT_STATUS" == ""]; then exit 0; else exit 1; fi