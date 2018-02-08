#!/bin/bash

NAME="vault-auth-file"
COMMIT=`git rev-parse --short HEAD`
TAG=`git describe --abbrev=0`
PACKAGE="github.com/gites/vault-auth-file/authfile"

LD_FLAGS="-X $PACKAGE.Name=$NAME -X $PACKAGE.GitCommit=$COMMIT -X $PACKAGE.Version=$TAG"
go build -ldflags "-s -w ${LD_FLAGS}"
