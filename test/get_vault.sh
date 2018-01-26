#!/bin/bash

VER=0.9.1
wget -q https://releases.hashicorp.com/vault/$VER/vault_${VER}_linux_amd64.zip
unzip vault_${VER}_linux_amd64.zip
chmod +x vault
