# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://vagrantcloud.com/search.
  config.vm.box = "bento/centos-7.4"

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  # config.vm.box_check_update = false

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # accessing "localhost:8080" will access port 80 on the guest machine.
  # NOTE: This will enable public access to the opened port
  # config.vm.network "forwarded_port", guest: 80, host: 8080

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine and only allow access
  # via 127.0.0.1 to disable public access
  #config.vm.network "forwarded_port", guest: 22, host: 22000, host_ip: "127.0.0.1"

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  # config.vm.network "private_network", ip: "192.168.33.10"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  # config.vm.synced_folder "../data", "/vagrant_data"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  # config.vm.provider "virtualbox" do |vb|
  #   # Display the VirtualBox GUI when booting the machine
  #   vb.gui = true
  #
  #   # Customize the amount of memory on the VM:
  #   vb.memory = "1024"
  # end
  #
  # View the documentation for the provider you are using for more
  # information on available options.

  # Enable provisioning with a shell script. Additional provisioners such as
  # Puppet, Chef, Ansible, Salt, and Docker are also available. Please see the
  # documentation for more information about their specific syntax and use.
   config.vm.provision "shell", inline: <<-SHELL
userdel bypass
userdel bypass-test1
userdel bypass-test2
userdel bypass-test3
userdel bypass-test4
userdel bypass-test5
userdel bypass-test6
userdel bypass-test7
userdel bypass-test8
userdel bypass-test9
userdel bypass-test10
rm -rf /home/bypass
rm -rf /home/bypass-test1
rm -rf /home/bypass-test2
rm -rf /home/bypass-test3
rm -rf /home/bypass-test4
rm -rf /home/bypass-test5
rm -rf /home/bypass-test6
rm -rf /home/bypass-test7
rm -rf /home/bypass-test8
rm -rf /home/bypass-test9
rm -rf /home/bypass-test10
useradd bypass 
useradd bypass-test1
useradd bypass-test2
useradd bypass-test3
useradd bypass-test4
useradd bypass-test5
useradd bypass-test6
useradd bypass-test7
useradd bypass-test8
useradd bypass-test9
useradd bypass-test10
echo 'fakepass' | passwd --stdin bypass
echo 'bypass ALL=(ALL:ALL) /bin/su - bypass-test1' | sudo EDITOR='tee -a' visudo
echo 'bypass ALL=(ALL:ALL) /bin/su - bypass-test2' | sudo EDITOR='tee -a' visudo
echo 'bypass ALL=(ALL:ALL) /bin/su - bypass-test3' | sudo EDITOR='tee -a' visudo
echo 'bypass ALL=(ALL:ALL) /bin/su - bypass-test4' | sudo EDITOR='tee -a' visudo
echo 'bypass ALL=(ALL:ALL) /bin/su - bypass-test5' | sudo EDITOR='tee -a' visudo
echo 'bypass ALL=(ALL:ALL) /bin/su - bypass-test6' | sudo EDITOR='tee -a' visudo
echo 'bypass ALL=(ALL:ALL) /bin/su - bypass-test7' | sudo EDITOR='tee -a' visudo
echo 'bypass ALL=(ALL:ALL) /bin/su - bypass-test8' | sudo EDITOR='tee -a' visudo
echo 'bypass ALL=(ALL:ALL) /bin/su - bypass-test9' | sudo EDITOR='tee -a' visudo
echo 'bypass ALL=(ALL:ALL) /bin/su - bypass-test10' | sudo EDITOR='tee -a' visudo
   SHELL
end
