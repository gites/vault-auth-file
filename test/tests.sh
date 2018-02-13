#!/bin/bash

GOPATH="${GOPATH:-$HOME/go}"

TEST_DIR="${GOPATH}/src/github.com/gites/vault-auth-file/test"

echo "GOPATH => $GOPATH"
echo "TEST_DIR -> $TEST_DIR"

set -ex

cd $TEST_DIR

go test -v ../authfile -timeout=60s -race -parallel=20 -covermode=atomic -coverprofile=coverage.txt

if [ ! -x vault ]; then
  ./get_vault.sh
fi

echo "plugin_directory=\"$GOPATH/src/github.com/gites/vault-auth-file\"" > config.hcl
./vault server -log-level=trace -dev -config config.hcl 2>stderr.log 1>stdout.log &
PID=$!
sleep 1
export VAULT_ADDR=http://127.0.0.1:8200
SHA_256SUM=`sha256sum ../vault-auth-file|cut -d' ' -f1`
./vault write sys/plugins/catalog/vault-auth-file \
        sha_256=$SHA_256SUM \
        command=vault-auth-file

./vault auth enable -path=file -plugin-name vault-auth-file plugin

./vault auth list

./vault audit enable file file_path=./vault_audit.log log_raw=true

set +e
./vault write auth/file/login username=wac password=lubieplacki && exit 1
set -e

./vault write auth/file/config path="$TEST_DIR/password-file"

./vault read auth/file/config

./vault write auth/file/config path="$TEST_DIR/password-file"

./vault read auth/file/config
# This one should fail
set +e
./vault write -format=json auth/file/login username=wac password=nielubieplackow && exit 1
set -e
./vault write -format=json auth/file/login username=wac password=lubieplacki

./vault login -method=userpass -path=file username=wac password=lubieplacki

./vault token renew

set +e
#cat *log
#sleep 3600
kill $PID
