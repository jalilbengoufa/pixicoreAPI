#!/bin/bash
set -eux

echo 'Defaults env_keep += "DEBIAN_FRONTEND"' >/etc/sudoers.d/env_keep_apt
chmod 440 /etc/sudoers.d/env_keep_apt
export DEBIAN_FRONTEND=noninteractive
apt-get update
#apt-get upgrade -y


# install vim.

apt-get install -y --no-install-recommends vim
cat >/etc/vim/vimrc.local <<'EOF'
syntax on
set background=dark
set esckeys
set ruler
set laststatus=2
set nobackup
EOF

# install golang 

wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz
tar -xvf go1.10.3.linux-amd64.tar.gz
sudo mv go /usr/local
export PATH=$PATH:/usr/local/go/bin

# install git

sudo apt install git -y

# install pixiecore
go get -v go.universe.tf/netboot/cmd/pixiecore

# run pixiecore client 

#sudo ~/go/bin/pixiecore api http://localhost:3000 --dhcp-no-bind

# run pixiecore api

cp /vagrant/pixicoreAPI ~/

# provision the DHCP server.
# see http://www.syslinux.org/wiki/index.php?title=PXELINUX

apt-get install -y --no-install-recommends isc-dhcp-server
cat>/etc/dhcp/dhcpd.conf<<'EOF'

default-lease-time 300;
max-lease-time 300;
option domain-name-servers 8.8.8.8, 8.8.4.4;
option subnet-mask 255.255.255.0;
option routers 10.1.1.1;
subnet 10.1.1.0 netmask 255.255.255.0 {
  range 10.1.1.100 10.1.1.254;
}

EOF
sed -i -E 's,^(INTERFACESv4=).*,\1"eth1",' /etc/default/isc-dhcp-server
sed -i -E  's/^INTERFACESv6=/#&/' /etc/default/isc-dhcp-server
cat>/usr/local/sbin/dhcp-event<<'EOF'
#!/bin/bash
# this is called when a lease changes state.
# NB you can see these log entries with journalctl -t dhcp-event
logger -t dhcp-event "argv: $*"
for e in $(env); do
  logger -t dhcp-event "env: $e"
done
EOF
chmod +x /usr/local/sbin/dhcp-event
systemctl restart isc-dhcp-server


# setup NAT.
# see https://help.ubuntu.com/community/IptablesHowTo

apt-get install -y iptables

# enable IPv4 forwarding.
sysctl net.ipv4.ip_forward=1
sed -i -E 's,^\s*#?\s*(net.ipv4.ip_forward=).+,\11,g' /etc/sysctl.conf

# NAT through eth0.
# NB use something like -s 10.1.1/24 to limit to a specific network.
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE

# load iptables rules on boot.
iptables-save >/etc/iptables-rules-v4.conf
cat<<'EOF'>/etc/network/if-pre-up.d/iptables-restore
#!/bin/sh
iptables-restore </etc/iptables-rules-v4.conf
EOF
chmod +x /etc/network/if-pre-up.d/iptables-restore


