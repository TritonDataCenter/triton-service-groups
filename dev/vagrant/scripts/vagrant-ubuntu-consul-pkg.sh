#!/bin/sh

set -o errexit

cd /tmp/
curl -sSL https://releases.hashicorp.com/consul/1.0.0/consul_1.0.0_linux_amd64.zip > consul.zip

echo "Installing Consul..."
unzip /tmp/consul.zip
sudo install consul /usr/bin/consul
(
cat <<-EOF
	[Unit]
	Description=consul agent
	Requires=network-online.target
	After=network-online.target

	[Service]
	Restart=on-failure
	ExecStart=/usr/bin/consul agent -dev
	ExecReload=/bin/kill -HUP $MAINPID

	[Install]
	WantedBy=multi-user.target
EOF
) | sudo tee /etc/systemd/system/consul.service

sudo systemctl enable consul.service
sudo systemctl start consul

sudo systemctl daemon-reload

for bin in cfssl cfssl-certinfo cfssljson
do
	echo "Installing $bin..."
	curl -sSL https://pkg.cfssl.org/R1.2/${bin}_linux-amd64 > /tmp/${bin}
	sudo install /tmp/${bin} /usr/local/bin/${bin}
done
