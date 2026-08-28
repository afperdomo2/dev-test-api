#!/bin/sh
test -z "$(gofmt -l .)" || { echo "Unformatted files:"; gofmt -l .; exit 1; }
