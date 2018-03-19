#!/bin/sh

set -o errexit

cd /tmp/
curl -sSL https://releases.hashicorp.com/nomad/0.7.1/nomad_0.7.1_linux_amd64.zip -o nomad.zip

unzip nomad.zip
sudo install nomad /usr/bin/nomad
sudo mkdir -p /etc/nomad.d
sudo chmod a+w /etc/nomad.d

sudo mkdir /etc/nomad
