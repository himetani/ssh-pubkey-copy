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
userdel test1
userdel common1
userdel common2
userdel common3
userdel common4
userdel common5
userdel common6
userdel common7
userdel common8
userdel common9
userdel common10
rm -rf /home/test1
rm -rf /home/common1
rm -rf /home/common2
rm -rf /home/common3
rm -rf /home/common4
rm -rf /home/common5
rm -rf /home/common6
rm -rf /home/common7
rm -rf /home/common8
rm -rf /home/common9
rm -rf /home/common10
useradd test1
useradd common1
useradd common2
useradd common3
useradd common4
useradd common5
useradd common6
useradd common7
useradd common8
useradd common9
useradd common10
echo 'fakepass' | passwd --stdin test1
echo 'test1 ALL=(ALL:ALL) /bin/su - common1' | sudo EDITOR='tee -a' visudo
echo 'test1 ALL=(ALL:ALL) /bin/su - common2' | sudo EDITOR='tee -a' visudo
echo 'test1 ALL=(ALL:ALL) /bin/su - common3' | sudo EDITOR='tee -a' visudo
echo 'test1 ALL=(ALL:ALL) /bin/su - common4' | sudo EDITOR='tee -a' visudo
echo 'test1 ALL=(ALL:ALL) /bin/su - common5' | sudo EDITOR='tee -a' visudo
echo 'test1 ALL=(ALL:ALL) /bin/su - common6' | sudo EDITOR='tee -a' visudo
echo 'test1 ALL=(ALL:ALL) /bin/su - common7' | sudo EDITOR='tee -a' visudo
echo 'test1 ALL=(ALL:ALL) /bin/su - common8' | sudo EDITOR='tee -a' visudo
echo 'test1 ALL=(ALL:ALL) /bin/su - common9' | sudo EDITOR='tee -a' visudo
echo 'test1 ALL=(ALL:ALL) /bin/su - common10' | sudo EDITOR='tee -a' visudo
   SHELL
end
