#!/bin/bash

ARG1=$1
VER=${ARG1:-0.9.3}
echo "Downloading v$VER Vault binary"
wget -q https://releases.hashicorp.com/vault/$VER/vault_${VER}_linux_amd64.zip
unzip vault_${VER}_linux_amd64.zip
chmod +x vault
