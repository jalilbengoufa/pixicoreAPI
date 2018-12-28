

Vagrant.configure("2") do |config|

  config.vm.box = "debian/stretch64"

  config.vm.provider :virtualbox do |vb|
    vb.linked_clone = true
    vb.memory = 256
    vb.customize ['modifyvm', :id, '--cableconnected2', 'on']
  end

  config.vm.define :master do |config|
    config.vm.hostname = 'master'
    config.vm.network "private_network", mac: '080027000011',ip: '10.1.1.1'
    config.vm.provision :shell, path: 'gateway.sh'
    config.vm.provider :virtualbox do |vb, config|
      vb.memory = 2048
      vb.customize ['modifyvm', :id, '--macaddress1', '080027000012']
	    vb.customize ['modifyvm', :id, '--natnet1', '10.0.2.0/24']
    end
  end

  config.vm.define :child do |config|
    config.vm.hostname = 'child'
    config.vm.network "private_network", mac: '080027000021',ip: '10.1.1.2', auto_config: false
    config.vm.provision "shell", inline: "sudo ip link set dev eth1 up && sudo dhclient eth1"
    config.vm.provider :virtualbox do |vb, config|
      vb.memory = 2048
      vb.customize ['modifyvm', :id, '--macaddress1', '080027000022']
	    vb.customize ['modifyvm', :id, '--natnet1', '10.0.2.0/24']
    end
  end

end 