#!/bin/sh

set -o errexit

chown vagrant:wheel \
       /opt/gopath \
       /opt/gopath/src \
       /opt/gopath/src/github.com \
       /opt/gopath/src/github.com/sean-

export ASSUME_ALWAYS_YES=yes

pkg update
pkg install -y \
	databases/cockroach \
	devel/git \
	editors/vim-console \
	lang/go \
	security/ca_root_nss \
	shells/bash \
	sysutils/tmux \
	sysutils/tree

chsh -s /usr/local/bin/bash vagrant
chsh -s /usr/local/bin/bash root

cat <<EOT >> /home/vagrant/.profile
export GOPATH=/opt/gopath
export PATH=\$GOPATH/bin:\$PATH

cd /opt/gopath/src/github.com/sean-/vpc
EOT

