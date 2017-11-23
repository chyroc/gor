#!/usr/bin/env bash

FMT_STATUS=`gofmt -d $(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./_project/*")`
if ["$FMT_STATUS" == ""]; then exit 0; else exit 1; fi
