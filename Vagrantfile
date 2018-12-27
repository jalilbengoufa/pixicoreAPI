# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|

  config.vm.box = "debian/stretch64"

  config.vm.provider :virtualbox do |vb|
    vb.linked_clone = true
    vb.memory = 256
    vb.customize ['modifyvm', :id, '--cableconnected1', 'on']
  end

  config.vm.define :gateway do |config|
    config.vm.hostname = 'gateway'
    config.vm.network :private_network, ip: '10.10.10.2'
    config.vm.provision :shell, path: 'gateway.sh'
  end

  config.vm.define :debian_live do |config|
    config.vm.box = 'debian/stretch64'
    config.vm.network :private_network, mac: '080027000001', ip: '10.10.10.3', auto_config: false
    config.vm.provider :virtualbox do |vb, config|
      # make sure this vm has enough memory to load the root fs into memory.
      vb.memory = 2048

      # let vagrant known that the guest does not have the guest additions nor a functional vboxsf or shared folders.
      vb.check_guest_additions = false
      vb.functional_vboxsf = false
      config.vm.synced_folder '.', '/vagrant', disabled: true

      # configure for PXE boot.
      vb.customize ['modifyvm', :id, '--boot1', 'net']
      vb.customize ['modifyvm', :id, '--boot2', 'disk']
      vb.customize ['modifyvm', :id, '--biospxedebug', 'on']
      vb.customize ['modifyvm', :id, '--cableconnected2', 'on']
      vb.customize ['modifyvm', :id, '--nicbootprio2', '1']
      vb.customize ['modifyvm', :id, "--nictype2", '82540EM'] # Must be an Intel card (as-of VB 5.1 we cannot Intel PXE boot from a virtio-net card).

      #
      # NB the VirtualBox BIOS is based on Plex86/Boch/QEMU.
      # NB dump extradata with VBoxManage getextradata $(cat .vagrant/machines/debianlive/virtualbox/id)
      vb.customize ['setextradata', :id, 'VBoxInternal/Devices/pcbios/0/Config/DmiSystemVendor',    'your vendor name here']
      vb.customize ['setextradata', :id, 'VBoxInternal/Devices/pcbios/0/Config/DmiSystemProduct',   'your product name here']
      vb.customize ['setextradata', :id, 'VBoxInternal/Devices/pcbios/0/Config/DmiSystemVersion',   'your product version here']
      vb.customize ['setextradata', :id, 'VBoxInternal/Devices/pcbios/0/Config/DmiSystemSerial',    'your product serial number here']
      vb.customize ['setextradata', :id, 'VBoxInternal/Devices/pcbios/0/Config/DmiSystemSKU',       'your product SKU here']
      vb.customize ['setextradata', :id, 'VBoxInternal/Devices/pcbios/0/Config/DmiSystemUuid',      '00000000-0000-4000-8000-000000000001']
      vb.customize ['setextradata', :id, 'VBoxInternal/Devices/pcbios/0/Config/DmiChassisVendor',   'your chassis vendor name here']
      vb.customize ['setextradata', :id, 'VBoxInternal/Devices/pcbios/0/Config/DmiChassisType',     '1']
      vb.customize ['setextradata', :id, 'VBoxInternal/Devices/pcbios/0/Config/DmiChassisVersion',  'your chassis version here']
      vb.customize ['setextradata', :id, 'VBoxInternal/Devices/pcbios/0/Config/DmiChassisSerial',   'your chassis serial number here']
      vb.customize ['setextradata', :id, 'VBoxInternal/Devices/pcbios/0/Config/DmiChassisAssetTag', 'your chassis asset tag here']
    end

    # make sure we use the gateway machine as this machine default gateway.
    # NB the 10.0.2/24 network is the default VirtualBox NAT network, which we must replace with our gateway.
    config.vm.provision :shell, run: 'always', inline: 'ip route list 0/0 | xargs ip route del; ip route add default via 10.10.10.2'

    # dump some useful information.
    config.vm.provision :shell, inline: '''
        set -x
        uname -a
        ip route
        ip route get 8.8.8.8
        cat /proc/cmdline
        cat /sys/devices/virtual/dmi/id/sys_vendor
        cat /sys/devices/virtual/dmi/id/product_name
        cat /sys/devices/virtual/dmi/id/product_version
        cat /sys/devices/virtual/dmi/id/product_serial
        cat /sys/devices/virtual/dmi/id/product_uuid
        cat /sys/devices/virtual/dmi/id/chassis_vendor
        cat /sys/devices/virtual/dmi/id/chassis_type
        cat /sys/devices/virtual/dmi/id/chassis_version
        cat /sys/devices/virtual/dmi/id/chassis_serial
        cat /sys/devices/virtual/dmi/id/chassis_asset_tag
      '''
end
  
end
