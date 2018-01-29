#!/bin/bash


set -ex
go test -v ../...

cd /home/wac/go/src/vault-auth-file/test
./vault server -dev -config config.hcl 2>stderr.log 1>stdout.log &
PID=$!
sleep 1
export VAULT_ADDR=http://127.0.0.1:8200
SHA_256SUM=`sha256sum /home/wac/go/src/vault-auth-file/vault-auth-file|cut -d' ' -f1`
./vault write sys/plugins/catalog/vault-auth-file \
        sha_256=$SHA_256SUM \
        command=vault-auth-file

./vault auth enable -path=file -plugin-name vault-auth-file plugin

./vault auth list

./vault audit enable file file_path=./vault_audit.log log_raw=true

set +e
./vault write auth/file/login username=wac password=lubieplacki && exit 1
set -e

./vault write auth/file/config path=/home/wac/go/src/vault-auth-file/test/password-file

./vault read auth/file/config

./vault write auth/file/config path=/home/wac/go/src/vault-auth-file/test/password-file ttl=123 max_ttl=456

./vault read auth/file/config
set +e
./vault write -format=json auth/file/login username=wac password=nielubueplackow && exit 1
set -e
./vault write -format=json auth/file/login username=wac password=lubieplacki | tee wac.json

./vault token renew `cat wac.json  | grep client_token | cut -d'"' -f 4`

set +e
#cat *log
#sleep 3600
kill $PID
rm wac.json
