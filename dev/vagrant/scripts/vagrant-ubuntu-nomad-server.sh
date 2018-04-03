#!/bin/sh

set -o errexit

(
cat <<-EOF
datacenter = "dc1"

# Increase log verbosity
log_level = "DEBUG"

# Setup data dir
data_dir = "/tmp/server1"

# Enable the server
server {
  enabled = true

  # Self-elect, should be 3 or 5 for production
  bootstrap_expect = 1
}

EOF
) | sudo tee /etc/nomad/server.hcl

(
cat <<-EOF
[Unit]
Description=Nomad
Documentation=https://nomadproject.io/docs/

[Service]
ExecStart=/usr/bin/nomad agent -config /etc/nomad/server.hcl
ExecReload=/bin/kill -HUP $MAINPID
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF
) | sudo tee /etc/systemd/system/nomad-server.service

sudo systemctl enable nomad-server.service
sudo systemctl start nomad-server.service

sudo systemctl daemon-reload
