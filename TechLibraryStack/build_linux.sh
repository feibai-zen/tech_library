#!/usr/bin/env bash

##########
VERSION=${1}

# These are the values we want to pass for Version and BuildTime
GITHASH=`git rev-parse HEAD 2>/dev/null`

BUILDAT=`date +%FT%T%z`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS="-s -w -X github.com/TruthHun/BookStack/utils.GitHash=${GITHASH} -X github.com/TruthHun/BookStack/utils.BuildAt=${BUILDAT} -X github.com/TruthHun/BookStack/utils.Version=${VERSION}"

##########

rm -rf output/${VERSION}
mkdir -p output/${VERSION}

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o output/${VERSION}/linux/BookStack -ldflags "${LDFLAGS}"

upx -f -9 output/${VERSION}/linux/BookStack

cp -r conf output/${VERSION}/linux/

cp -r views output/${VERSION}/linux/

cp -r dictionary output/${VERSION}/linux/

cp -r static output/${VERSION}/linux/
