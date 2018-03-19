#!/bin/sh

set -o errexit

(
cat <<-EOF
# Increase log verbosity
log_level = "DEBUG"

# Setup data dir
data_dir = "/tmp/client1"

# Enable the client
client {
  enabled = true

  # For demo assume we are talking to server1. For production,
  # this should be like "nomad.service.consul:4647" and a system
  # like Consul used for service discovery.
  servers = ["127.0.0.1:4647"]
}

# Modify our port to avoid a collision with server1
ports {
  http = 5656
}

EOF
) | sudo tee /etc/nomad/client.hcl


(
cat <<-EOF
[Unit]
Description=Nomad
Documentation=https://nomadproject.io/docs/

[Service]
ExecStart=/usr/bin/nomad agent -config /etc/nomad/client.hcl
ExecReload=/bin/kill -HUP $MAINPID
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF
) | sudo tee /etc/systemd/system/nomad-client.service

sudo systemctl enable nomad-client.service
sudo systemctl start nomad-client.service

sudo systemctl daemon-reload
