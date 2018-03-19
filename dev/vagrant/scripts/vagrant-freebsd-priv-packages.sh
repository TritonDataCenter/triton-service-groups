#!/bin/sh

set -o errexit

export ASSUME_ALWAYS_YES=yes

pkg update
pkg install -y \
	dns/nss_mdns \
	editors/vim-console \
	net/avahi-app \
	security/ca_root_nss \
	shells/bash \
	sysutils/grub2-bhyve \
	sysutils/tmux \
	sysutils/tree

chsh -s /usr/local/bin/bash vagrant
chsh -s /usr/local/bin/bash root
