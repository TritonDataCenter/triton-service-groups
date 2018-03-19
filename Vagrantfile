# -*- mode: ruby -*-
# vi: set ft=ruby :

crdbNodes = [
    { :name => "crdb1", :ip => "172.27.10.11", :autostart => true, },
    { :name => "crdb2", :ip => "172.27.10.12", :autostart => true, },
    { :name => "crdb3", :ip => "172.27.10.13", :autostart => true, },
]

nomadNodes = [
   { :name => "nomad", :ip => "172.27.10.21", :autostart => true, },
]

freebsd_box = 'jen20/FreeBSD-12.0-CURRENT-VPC'
ubuntu_box = 'bento/ubuntu-16.04'

require './dev/vagrant/helper/core'
require './dev/vagrant/helper/utils'

Vagrant.configure("2") do |config|
	config.ssh.extra_args = ["-e", "%"]

	crdbNodes.each do |node|
	    config.vm.define node[:name], autostart: node[:autostart] do |vmCfg|
            vmCfg.vm.box = freebsd_box
			      vmCfg.vm.hostname = node[:name]
			      vmCfg = configureFreeBSDDBProvisioners(vmCfg, node[:name], node[:ip])

            vmCfg = addPrivateNICOptions(vmCfg, node[:ip])
            vmCfg = configureMachineSize(vmCfg, 2, 1024)
	    end
	end

	nomadNodes.each do |node|
        config.vm.define node[:name], autostart: node[:autostart] do |vmCfg|
            vmCfg.vm.box = ubuntu_box
            vmCfg.vm.hostname = node[:name]

            vmCfg = configureUbuntuNomadProvisioners(vmCfg, node[:name], node[:ip])
            vmCfg = configureMachineSize(vmCfg, 2, 1024)

            vmCfg.vm.network "forwarded_port", guest: 4646, host: 4646, auto_correct: true
        end
    end
end

def addPrivateNICOptions(vmCfg, ip)
	vmCfg.vm.network "private_network", ip: ip

	["vmware_fusion", "vmware_workstation"].each do |p|
		vmCfg.vm.provider p do |v|
			v.vmx["ethernet1.virtualdev"] = "vmxnet3"
			v.vmx["ethernet1.pcislotnumber"] = "192"
		end
	end

	return vmCfg
end

def configureMachineSize(vmCfg, vcpuCount, memSize)
	["vmware_fusion", "vmware_workstation"].each do |p|
		vmCfg.vm.provider p do |v|
			v.vmx["memsize"] = "1024"
			v.vmx["numvcpus"] = "2"
		end
	end

	return vmCfg
end

def configureUbuntuNomadProvisioners(vmCfg, hostname, ip)
    vmCfg.vm.provision "shell",
    		path: './dev/vagrant/scripts/vagrant-ubuntu-deps.sh',
    		privileged: true

    vmCfg.vm.provision "shell",
        path: './dev/vagrant/scripts/vagrant-ubuntu-nomad-pkg.sh',
        privileged: true

    vmCfg.vm.provision "shell",
        path: './dev/vagrant/scripts/vagrant-ubuntu-consul-pkg.sh',
        privileged: true

    vmCfg.vm.provision "shell",
        path: './dev/vagrant/scripts/vagrant-ubuntu-nomad-server.sh',
        privileged: true

    vmCfg.vm.provision "shell",
        path: './dev/vagrant/scripts/vagrant-ubuntu-nomad-client.sh',
        privileged: true

    return vmCfg
end

def configureFreeBSDDBProvisioners(vmCfg, hostname, ip)
	vmCfg.vm.provision "shell",
		path: './dev/vagrant/scripts/vagrant-freebsd-priv-db-packages.sh',
		privileged: true

	vmCfg.vm.provision "file",
		source: './dev/vagrant/certs/ca/ca.crt',
		destination: "/home/vagrant/.cockroach-certs/ca.crt"

	vmCfg.vm.provision "file",
		source: "./dev/vagrant/certs/client/client.root.crt",
		destination: "/home/vagrant/.cockroach-certs/client.root.crt"

	vmCfg.vm.provision "file",
		source: "./dev/vagrant/certs/client/client.root.key",
		destination: "/home/vagrant/.cockroach-certs/client.root.key"

	vmCfg.vm.provision "file",
		source: './dev/vagrant/certs/ca/ca.crt',
		destination: "/secrets/crdb/ca.crt"

	vmCfg.vm.provision "file",
		source: "./dev/vagrant/certs/#{hostname}/node.crt",
		destination: "/secrets/crdb/node.crt"

	vmCfg.vm.provision "file",
		source: "./dev/vagrant/certs/#{hostname}/node.key",
		destination: "/secrets/crdb/node.key"

	vmCfg.vm.provision "shell",
		path: './dev/vagrant/scripts/vagrant-freebsd-priv-db-configure.sh',
		privileged: true

	if hostname == "crdb3"
		vmCfg.vm.provision "shell",
			path: './dev/vagrant/scripts/vagrant-freebsd-unpriv-db-init.sh',
			privileged: false
	end

	return vmCfg
end
