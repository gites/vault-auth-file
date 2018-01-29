#!/bin/bash

ARG1=$1
VER=${ARG1:-0.9.2}

wget -q https://releases.hashicorp.com/vault/$VER/vault_${VER}_linux_amd64.zip
unzip vault_${VER}_linux_amd64.zip
chmod +x vault
